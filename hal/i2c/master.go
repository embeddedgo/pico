// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c

import (
	"embedded/rtos"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/embeddedgo/device/bus/i2cbus"
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/system/clock"
)

// A Master is a driver for the I2C peripheral. It provides two kinds of
// interfaces to communicate with slave devices on the I2C bus.
//
// The first interface is a low-level one. It provides a set of methods to
// directly interract with the Data / Command FIFOs of the underlying I2C
// peripheral.
//
// Example:
//
//	d.SetAddr(eepromAddr)
//	d.WriteCmds([]int16{
//		lpi2c.Send|int16(memAddr),
//		lpi2c.Recv|int16(len(buf) - 1),
//		lpi2c.Stop,
//	})
//	d.ReadBytes(buf)
//	if err := d.Err(true); err != nil {
//
// Write methods in the low-level interface are asynchronous, that is, they may
// return before all commands/data will be written to the FIFO. Therefore you
// must not modify the data/command buffer passed to the last write method until
// the return of the Flush method or another write method.
//
// The read/write methods doesn't return errors. There is an Err method that
// allow to check and reset the I2C error flags at a convenient time. Even if
// you call Err after every method call the returned error is still asynchronous
// due to the asynchronous nature of the write methods and the delayed execution
// of commands by the I2C peripheral itself. You can use Wait before checking
// error (especially waiting for STOP_DET) to somehow synhronize things.
//
// The second interface is a connection oriented one that implements the
// i2cbus.Conn interface.
//
// Example:
//
//	c := d.NewConn(eepromAddr)
//	c.WriteByte(memAaddr)
//	c.Read(buf)
//	err := c.Close()
//	if err != nil {
//
// Both interfaces may be used concurently by multiple goroutines but in such
// a case users of the low-level interface must gain an exclusive access to the
// driver using the embedded mutex and wait for the Stop Condition before
// unlocking the Master.
type Master struct {
	sync.Mutex

	name string
	p    *Periph

	wbuf int16
	id   uint8

	cmd   bool
	wdata unsafe.Pointer
	wi    int32 // ISR cannot alter the above pointers so it alters wi instead
	wn    int32
	wdone rtos.Note

	rdata *byte
	ri    int32 // ISR cannot alter the above pointer so it alters ri instead
	rn    int32
	rdone rtos.Note

	dma dma.Channel
	dcf dma.Config
	din int
}

// NewMaster returns a new master-mode driver for p. If valid DMA channel is
// given, the DMA will be used for bigger data transfers.
func NewMaster(p *Periph, dc dma.Channel) *Master {
	req := dma.I2C0_TX + dma.Config(num(p))*(dma.I2C1_TX-dma.I2C0_TX)
	return &Master{
		name: "I2C" + string(rune('0'+num(p))),
		p:    p,
		dma:  dc,
		dcf:  dma.En | dma.S8b | req,
		din:  int(system.NextCPU()),
	}
}

// Periph returns the underlying SPI peripheral.
func (d *Master) Periph() *Periph {
	return d.p
}

// Setup resets and configures the underlying I2C pripheral to operate in the
// master mode with the given speed.
func (d *Master) Setup(baudrate int) {
	p := d.p
	p.SetReset(true)
	p.SetReset(false)
	p.INTR_MASK.Store(0)

	// Always use FAST mode as PICO-SDK does.
	p.CON.Store(MASTER_MODE | SLAVE_DISABLE | RESTART_EN | TX_EMPTY_CTRL | RX_FIFO_FULL_HLD_CTRL | FAST)
	p.DMA_CR.Store(TDMAE | RDMAE) // enable by defalut, on/off at the DMA side

	// Baudrate (calculations taken from PICO-SDK)

	clk := clock.PERI.Freq()
	cn := uint32((clk + int64(baudrate/2)) / int64(baudrate))

	lcn := cn * 3 / 5
	hcn := cn - lcn

	var txHold SDA_HOLD
	if baudrate < 1e6 {
		txHold = SDA_HOLD(clk*3/10e6) + 1
	} else {
		txHold = SDA_HOLD(clk*3/25e6) + 1
	}

	spkn := uint32(1)
	if lcn >= 16 {
		spkn = lcn / 16
	}

	p.FS_SCL_LCNT.Store(lcn)
	p.FS_SCL_HCNT.Store(hcn)
	p.FS_SPKLEN.Store(spkn)
	p.SDA_HOLD.StoreBits(SDA_TX_HOLD, txHold)
}

// SetAddr sets the address of a slave device. You must ensure there is no
// any command in the Tx FIFO intended to use the previous address (a command
// that causes Start or Repeated Start Condition).
func (d *Master) SetAddr(addr i2cbus.Addr) {
	p := d.p
	p.ENABLE.Store(0)
	if addr&i2cbus.A10 != 0 {
		addr &= 0x3ff
		internal.AtomicSet(&p.CON, MASTER_10BITADDR)
	} else {
		addr &= 0x7f
		internal.AtomicClear(&p.CON, MASTER_10BITADDR)
	}
	p.TAR.Store(TAR(addr))
	p.ENABLE.Store(EN)
}

const (
	Send    = int16(0)
	Recv    = int16(CMD)
	Stop    = int16(STOP)
	Restart = int16(RESTART)
)

const (
	rxFIFOCap = 16
	txFIFOCap = 16
	abrtFlags = 0x1ffff
	minDMA    = 16
)

// Flush waits until all commands/data passed to the driver have been consumed
// (in other words, it makes the previous write operation synchronous). You must
// call Flush or write new to enusre the Master stops referencing previously
// written data (to reuse memory or make it available for garbage collection).
// Return from Flush doesn't mean that all data were sent on the bus (there may
// be even full Tx FIFO not handled yet, see Wait).
func (d *Master) Flush() {
	if d.wdata != nil {
		d.wdone.Sleep(-1)
		d.wdone.Clear()
		d.wdata = nil
	}
}

func masterWrite(d *Master, ptr unsafe.Pointer, n int, cmd bool) {
	p := d.p
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) != 0 {
		return
	}
	// To speed things up we try to fill the FIFO in thread mode. As thread code
	// may be interrupted at any time we check the TFNF bit every iteration
	// instead of write as fast as possible the txFIFOCap-TXFLR commands/bytes.
	i := 0
	if !cmd {
		data := unsafe.Slice((*byte)(ptr), n)
		for p.STATUS.LoadBits(TFNF) != 0 {
			p.DATA_CMD.Store(uint32(data[i]))
			if i++; i == len(data) {
				return
			}
		}
	} else {
		cmds := unsafe.Slice((*int16)(ptr), n)
		for ; ; i++ {
			if i >= len(cmds) {
				return
			}
			if p.STATUS.LoadBits(TFNF) == 0 {
				break
			}
			cmd := uint32(cmds[i])
			p.DATA_CMD.Store(cmd)
			if cmd&CMD == 0 {
				continue
			}
			// Handle the recvie length encoded as cmd&255 + 1.
		multiRecv:
			if cmd--; cmd&CMD == 0 {
				continue
			}
			if p.STATUS.LoadBits(TFNF) == 0 {
				cmds[i] = int16(cmd)
				break
			}
			p.DATA_CMD.Store(cmd)
			goto multiRecv
		}
	}
	// The remaining data/commands will be writtend to the FIFO by the ISR.
	d.cmd = cmd
	d.wdata = ptr
	d.wi = int32(i)
	atomic.StoreInt32(&d.wn, int32(n))
	internal.AtomicSet(&p.INTR_MASK, TX_EMPTY|TX_ABRT) // race with ISR (clear)
}

func masterWriteDMA(d *Master, ptr unsafe.Pointer, n int) {
	p := d.p
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) != 0 {
		return
	}
	d.wdata = ptr                // prevent premature GC and make Flush working
	atomic.StoreInt32(&d.wn, -1) // DMA write in progress
	dc := d.dma
	dc.ClearIRQ()
	dc.SetReadAddr(ptr)
	dc.SetWriteAddr(unsafe.Pointer(p.DATA_CMD.Addr() | 1<<14))
	dc.SetTransCount(n, dma.Normal)
	dc.SetConfigTrig(d.dcf|dma.IncR, dc)
	dc.EnableIRQ(d.din)
	internal.AtomicSet(&p.INTR_MASK, TX_ABRT)
}

// WriteCmd works like WriteCmds but writes only one command word into the Tx
// FIFO.
func (d *Master) WriteCmd(cmd int16) {
	d.Flush()
	d.wbuf = cmd
	masterWrite(d, unsafe.Pointer(&d.wbuf), 1, true)
}

// WriteCmds starts writing commands into the Tx FIFO in the background using
// interrupts and/or DMA. WriteCmd is no-op if len(cmds) == 0.
func (d *Master) WriteCmds(cmds []int16) {
	if len(cmds) == 0 {
		return
	}
	d.Flush()
	masterWrite(d, unsafe.Pointer(unsafe.SliceData(cmds)), len(cmds), true)
}

// WriteBytes is like WriteCmds but writes only Send commands with the provided
// data.
func (d *Master) WriteBytes(p []byte) {
	if len(p) == 0 {
		return
	}
	d.Flush()
	if len(p) < minDMA || !d.dma.IsValid() {
		masterWrite(d, unsafe.Pointer(&p[0]), len(p), false)
	} else {
		masterWriteDMA(d, unsafe.Pointer(&p[0]), len(p))
	}
}

// WriteStr is like WriteBytes but writes bytes from string instead of slice.
func (d *Master) WriteStr(s string) {
	if len(s) != 0 {
		d.WriteBytes(unsafe.Slice(unsafe.StringData(s), len(s)))
	}
}

func masterRead(d *Master, ptr *byte, n int) {
	p := d.p
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) != 0 {
		return
	}
	// To speed things up we try to empty the FIFO in thread mode. As thread
	// code may be interrupted at any time we check the RFNE bit every iteration
	// instead of read as fast as possible the RXFLR bytes.
	i := 0
	data := unsafe.Slice((*byte)(ptr), n)
	for p.STATUS.LoadBits(RFNE) != 0 {
		data[i] = byte(p.DATA_CMD.Load())
		if i++; i == len(data) {
			return
		}
	}
	// The remaining data will be read by the ISR.
	d.rdata = ptr
	d.ri = int32(i)
	p.RX_TL.Store(uint32(min(n-i, rxFIFOCap) - 1))
	atomic.StoreInt32(&d.rn, int32(n))
	internal.AtomicSet(&p.INTR_MASK, RX_FULL|TX_ABRT) // race with ISR (clear)
	d.rdone.Sleep(-1)
	d.rdone.Clear()
	d.rdata = nil
}

func masterReadDMA(d *Master, ptr *byte, n int) {
	// Avoid GC the ptr prematurely if the caller intention is to read and
	// discard data so it may no longer reference the ptr.
	_masterReadDMA(d, uintptr(unsafe.Pointer(ptr)), n)
}

//go:uintptrescapes
func _masterReadDMA(d *Master, ptr uintptr, n int) {
	if atomic.LoadInt32(&d.wn) == -1 {
		d.Flush() // wait for the end of DMA write
	}
	p := d.p
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) != 0 {
		return
	}
	atomic.StoreInt32(&d.rn, -1) // DMA read in progress
	dc := d.dma
	dc.ClearIRQ()
	dc.SetReadAddr(unsafe.Pointer(p.DATA_CMD.Addr()))
	dc.SetWriteAddr(unsafe.Pointer(ptr))
	dc.SetTransCount(n, dma.Normal)
	dc.SetConfigTrig((d.dcf|dma.IncW)+(dma.I2C0_RX-dma.I2C0_TX), dc)
	dc.EnableIRQ(d.din)
	internal.AtomicSet(&p.INTR_MASK, TX_ABRT)
	d.rdone.Sleep(-1) // wait until the major loop complete or error
	d.rdone.Clear()
}

// ReadBytes reads len(p) data bytes from Rx FIFO. The read data is valid if Err
// returns nil.
func (d *Master) ReadBytes(p []byte) {
	if len(p) == 0 {
		return
	}
	if len(p) < minDMA || !d.dma.IsValid() {
		masterRead(d, &p[0], len(p))
	} else {
		masterReadDMA(d, &p[0], len(p))
	}
}

// ReadByte works like ReadBytes but reads only one byte from the Rx FIFO.
func (d *Master) ReadByte() (b byte) {
	masterRead(d, &b, 1)
	return
}

const statusFlags = RX_FULL | TX_EMPTY | TX_ABRT | ACTIVITY | STOP_DET | START_DET

// Status returns the flags that correspond to the current I2C master state
// (RX_FULL, TX_EMPTY) and the registered events (TX_ABRT, ACTIVITY, STOP_DET,
// START_DET). It is intended to be used together with the Clear and Wait
// methods. See also the documentation of the RAW_INTR_STAT register.
func (d *Master) Status() INTR {
	return d.p.RAW_INTR_STAT.LoadBits(statusFlags)
}

// Clear allows to clear the registered events except the TX_ABRT that can be
// cleared using the Err method. See Status for more information.
func (d *Master) Clear(flags INTR) {
	p := d.p
	if flags&ACTIVITY != 0 {
		p.CLR_ACTIVITY.Load()
	}
	if flags&STOP_DET != 0 {
		p.CLR_STOP_DET.Load()
	}
	if flags&START_DET != 0 {
		p.CLR_START_DET.Load()
	}
}

// Wait waits for an event/state specified by flags. See Status for more
// information.
func (d *Master) Wait(flags INTR) {
	flags &= statusFlags
	if flags == 0 {
		return
	}
	p := d.p
	if p.RAW_INTR_STAT.LoadBits(flags) != 0 {
		return
	}
	atomic.StoreInt32(&d.rn, -int32(flags))
	internal.AtomicSet(&p.INTR_MASK, flags)
	d.rdone.Sleep(-1)
	d.rdone.Clear()
}

// Err returns the content of the TX_ABRT_SOURCE register wrapped into the
// MasterError type if any sbort flag is set. Othewrise it returns nil. If clear
// is true Err clears the TX_ABRT_SOURCE register.
func (d *Master) Err(clear bool) (err error) {
	p := d.p
	if abort := p.TX_ABRT_SOURCE.Load(); abort&abrtFlags != 0 {
		err = &MasterError{d.name, abort}
		if clear {
			p.CLR_TX_ABRT.Load()
		}
	}
	return
}

// Abort aborts the I2C transfer. It can be used togather with Wait(TX_EMPTY) to
// implement asynchronous Stop condition. The command set supports only
// synchronous Stop by setting the Stop bit in the last send/receive command
// (you need to know in advance which command is the last command in the I2C
// transaction which isn't always convenient/possible).
func (d *Master) Abort() {
	p := d.p
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) != 0 {
		return
	}
	internal.AtomicSet(&p.ENABLE, ABORT)
	runtime.Gosched() // takes enough time so ABRT_USER_ABRT is almost always set
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) == ABRT_USER_ABRT {
		p.CLR_TX_ABRT.Load()
		return
	}
	// More expensive waiting.
	atomic.StoreInt32(&d.rn, -int32(TX_ABRT))
	internal.AtomicSet(&p.INTR_MASK, TX_ABRT)
	d.rdone.Sleep(-1)
	d.rdone.Clear()
	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags) == ABRT_USER_ABRT {
		p.CLR_TX_ABRT.Load()
	}
	return
}

// ISR is the interrupt handler for the I2C peripheral used by Master.
//
//go:nosplit
//go:nowritebarrierrec
func (d *Master) ISR() {
	p := d.p

	// Disable interrupts and reenable them later if needed. It reaces with the
	// thread code. If the clearing INTR_MASK here happens before the setting
	// it in the thread code this ISR may run again. We clear d.wn, d.rn before
	// wake-up the thread code so such ISR reentry isn't harmful.
	p.INTR_MASK.Store(0)

	if p.TX_ABRT_SOURCE.LoadBits(abrtFlags&^ABRT_USER_ABRT) != 0 {
		// Tx/Rx FIFOs are kept empty until TX_ABRT IRQ is cleared
		wn := atomic.LoadInt32(&d.wn)
		rn := atomic.LoadInt32(&d.rn)
		if wn|rn == -1 {
			dc := d.dma
			dc.DisableIRQ(d.din)
			dc.Abort()
		}
		if wn != 0 {
			d.wn = 0
			d.wdone.Wakeup()
		}
		if rn != 0 {
			d.rn = 0
			d.rdone.Wakeup()
		}
		return
	}

	var enable INTR

	// Read or wait part.
	done := false
	if n := atomic.LoadInt32(&d.rn); n > 0 {
		// Read
		flags := RX_FULL | TX_ABRT
		if fr := p.RXFLR.Load(); fr != 0 {
			i := d.ri
			m := min(n, int32(fr)+i)
			data := unsafe.Slice(d.rdata, m)
			for int(i) < len(data) {
				data[i] = byte(p.DATA_CMD.Load())
				i++
			}
			d.ri = i
			n -= i
			if n == 0 {
				flags = 0
				done = true
			} else {
				if n < rxFIFOCap {
					// Reduce the Rx threshold to the size of the last chunk
					p.RX_TL.Store(uint32(n - 1))
				}
			}
		}
		enable |= flags
	} else if n < 0 {
		// Wait
		if flags := INTR(-n); p.RAW_INTR_STAT.LoadBits(flags) != 0 {
			done = true
		} else {
			enable |= flags
		}
	}
	if done {
		d.rn = 0
		d.rdone.Wakeup()
	}

	// Write part. May work concurently with the thread read code.
	if n := atomic.LoadInt32(&d.wn); n > 0 {
		flags := TX_EMPTY | TX_ABRT
		if fw := txFIFOCap - p.TXFLR.Load(); fw != 0 {
			i := d.wi
			m := min(n, int32(fw)+i)
			if !d.cmd {
				for _, b := range unsafe.Slice((*byte)(d.wdata), m)[i:] {
					p.DATA_CMD.Store(uint32(b))
				}
			} else {
				cmds := unsafe.Slice((*int16)(d.wdata), m)
				for ; i < m; i++ {
					cmd := uint32(cmds[i])
					p.DATA_CMD.Store(cmd)
					if cmd&CMD == 0 {
						continue
					}
					// Handle the recvie length encoded as cmd&255 + 1.
				multiRecv:
					if cmd--; cmd&CMD == 0 {
						continue
					}
					if m--; i == m {
						cmds[i] = int16(cmd)
						break
					}
					p.DATA_CMD.Store(cmd)
					goto multiRecv
				}
			}
			d.wi = m
			if m == n {
				// Done.
				flags = 0
				d.wn = 0
				d.wdone.Wakeup()
			}
		}
		enable |= flags
	}

	// Reenable interrupts for unfinished requests.
	if enable != 0 {
		internal.AtomicSet(&p.INTR_MASK, enable)
	}
}

// DMAISR is a DMA interrupt handler for the DMA channel used by Master.
//
//go:nosplit
//go:nowritebarrierrec
func (d *Master) DMAISR() {
	d.dma.DisableIRQ(d.din)
	if atomic.LoadInt32(&d.wn) == -1 {
		d.wn = 0
		d.wdone.Wakeup()
	} else if atomic.LoadInt32(&d.rn) == -1 {
		d.rn = 0
		d.rdone.Wakeup()
	}
}
