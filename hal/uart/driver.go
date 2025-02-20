// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"embedded/rtos"
	"errors"
	"time"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/hal/system/clock"
)

// Driver provides a driver for the UART peripheral.
//
// The driver Rx and Tx parts are independent in means that they can be used
// conncurently by two goroutines. However, a given transmission direction can
// only be used by one goroutine at the same time.
type Driver struct {
	p *Periph

	wtimeout time.Duration
	wstart   uintptr
	wend     uintptr
	wdone    rtos.Note

	rtimeout time.Duration
	rstart   uintptr
	rend     uintptr
	rerr     uint32
	rready   rtos.Note
}

// NewDriver returns a new driver for p.
func NewDriver(p *Periph) *Driver {
	return &Driver{p: p, wtimeout: -1, rtimeout: -1}
}

// Periph returns the underlying LPSPI peripheral.
func (d *Driver) Periph() *Periph {
	return d.p
}

func (d *Driver) Enable() {
	d.p.CR.SetBits(UARTEN)
}

func (d *Driver) Disable() {
	d.p.CR.ClearBits(UARTEN)
}

type Config uint32

const (
	Break      = Config(BRK)        // send break
	ParityEven = Config(PEN | EPS)  // even parity
	ParityOdd  = Config(PEN)        // odd parity
	Stop2b     = Config(STP2)       // two stop bits instead of one
	Word5b     = Config(2 << WLENn) // 5 bit data word
	Word6b     = Config(2 << WLENn) // 6 bit data word
	Word7b     = Config(2 << WLENn) // 7 bit data word
	Word8b     = Config(3 << WLENn) // 8 bit data word
)

func (d *Driver) Config() Config {
	return Config(d.p.LCR_H.Load())
}

// SetConfig sets the URAT configuration.
func (d *Driver) SetConfig(cfg Config) {
	p := d.p
	cr := p.CR.Load()
	p.CR.Store(cr &^ UARTEN) // disable UART before accessing LCR
	p.LCR_H.Store(LCR_H(cfg) | FEN)
	p.CR.Store(cr)
}

func (d *Driver) Baudrate() int {
	periHz := clock.PERI.Freq()
	p := d.p
	ibrd := p.IBRD.Load()
	fbrd := p.FBRD.Load()
	div := int64(ibrd<<6 + fbrd)
	return int((uint(8*periHz/div) + 1) >> 1)
}

// SetBaudrate sets the UART baudrate.
func (d *Driver) SetBaudrate(baudrate int) {
	periHz := clock.PERI.Freq()
	brdiv := uint32(8*periHz/int64(baudrate)) + 1
	ibrd := brdiv >> 7
	fbrd := brdiv & 0x7f >> 1
	if ibrd == 0 {
		ibrd = 1
		fbrd = 0
	} else if ibrd >= 0xffff {
		ibrd = 0xffff
		fbrd = 0
	}
	p := d.p
	cr := p.CR.Load()
	p.CR.Store(cr &^ UARTEN)      // disable UART before accessing LCR
	p.IBRD.Store(ibrd)            // IBRD is internally part of LCR
	p.FBRD.Store(fbrd)            // FBRD is internally part of LCR
	p.LCR_H.Store(p.LCR_H.Load()) // dummy write to latch IBRD and FBRD
	p.CR.Store(cr)
}

// Setup resets the underlying UART peripheral and configures it according to
// the driver needs. Next it calls the SetConfig and SetBaudrate methods with
// the provided arguments. You still need to call EnableTx/EnabeRx.
func (d *Driver) Setup(cfg Config, baudrate int) {
	p := d.p
	p.SetReset(true)
	p.SetReset(false)
	p.CR.Store(0) // disable Rx and Tx
	d.SetConfig(cfg)
	d.SetBaudrate(baudrate)
}

const fifoLen = 32

// ISR is the interrupt handler that handles the data transfers scheduled by
// the read and write methods.
//
//go:nosplit
//go:nowritebarrierrec
func (d *Driver) ISR() {
	p := d.p
	irqs := p.MIS.Load()
	p.ICR.Store(irqs) // clear active interupts

	// Write part.
	if irqs&TXI != 0 {
		// Write fast up to fifoLen/2 bytes to the FIFO.
		addr, end := d.wstart, d.wend
		fast := min(addr+fifoLen/2, end)
		for addr < fast {
			b := *(*byte)(unsafe.Pointer(addr))
			p.DR.Store(uint32(b))
			addr++
		}
		// Fill the FIFO to the end to reduce the interrupt rate.
		for addr < end && p.FR.LoadBits(TXFF) == 0 {
			b := *(*byte)(unsafe.Pointer(addr))
			p.DR.Store(uint32(b))
			addr++
		}
		d.wstart = addr
		if addr == end {
			// Transfer completed.
			internal.AtomicClear(&p.IMSC, TXI) // disable Tx FIFO interrupt
			d.wdone.Wakeup()
		}
	}

	// Read part
	if irqs&(RXI|RTI) != 0 {
		addr, end := d.rstart, d.rend
		if irqs&RXI != 0 {
			// Read fast up to fifoLen/2 bytes from the FIFO.
			fast := min(addr+fifoLen/2, end)
			for addr < fast {
				v := p.DR.Load()
				*(*byte)(unsafe.Pointer(addr)) = byte(v)
				if e := v >> 8 & 15; e != 0 {
					d.rerr = e
					goto end
				}
				addr++
			}
		}
		// Read the FIFO to the end to avoid overruns.
		for addr < end && p.FR.LoadBits(RXFE) == 0 {
			v := p.DR.Load()
			*(*byte)(unsafe.Pointer(addr)) = byte(v)
			if e := v >> 8 & 15; e != 0 {
				d.rerr = e
				goto end
			}
			addr++
		}
		d.rstart = addr
	end:
		internal.AtomicClear(&p.IMSC, RXI|RTI) // disable Rx FIFO interrupts
		d.rready.Wakeup()
	}
}

var ErrTimeout = errors.New("uart: timeout")

// SetReadTimeout sets the read timeout used by Read* functions.
func (d *Driver) SetReadTimeout(timeout time.Duration) {
	d.rtimeout = timeout
}

// SetWriteTimeout sets the write timeout used by Write* functions.
func (d *Driver) SetWriteTimeout(timeout time.Duration) {
	d.wtimeout = timeout
}

