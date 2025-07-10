// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"embedded/rtos"
	"runtime"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
)

// EnableTx enables the Tx part of the UART.
func (d *Driver) EnableTx() {
	internal.AtomicSet(&d.p.CR, UARTEN|TXE)
}

// DisableTx waits for the end of transfer (see WaitTxDone) and disables the Tx part
// of the UART.
func (d *Driver) DisableTx() {
	d.WaitTxDone()
	internal.AtomicClear(&d.p.CR, TXE)
}

//go:nosplit
func (d *Driver) Write(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	p := d.p
	// Because of the interrupt cost write FIFO in thread mode if possible.
	if len(s) > 3 {
		// The overhead required to setup a tight write loop may pay off.
		var m int
		if p.FR.LoadBits(TXFE) != 0 {
			m = fifoLen
		} else if p.RIS.LoadBits(TXI) != 0 {
			m = fifoLen / 2
		}
		if m != 0 {
			// There are at least m free locations in the FIFO.
			n = min(m, len(s))
			for _, b := range s[:n] {
				p.DR.Store(uint32(b))
			}
			if n == len(s) {
				return
			}
		}
	}
pollFF:
	for p.FR.LoadBits(TXFF) == 0 {
		// There is at least 1 free location in the FIFO.
		p.DR.Store(uint32(s[n]))
		if n++; n >= len(s) {
			return
		}
	}
	if rtos.HandlerMode() {
		// Write called in handler mode by print or println.
		goto pollFF
	}
	// The remaining data will be written to the FIFO by the ISR.
	if waitWriteISR(d, &s[n], len(s)-n) {
		return len(s), nil
	} else {
		return n, ErrTimeout
	}
}

func (d *Driver) WriteString(s string) (n int, err error) {
	return d.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (d *Driver) WriteByte(b byte) error {
	p := d.p
	if p.FR.LoadBits(TXFF) == 0 {
		p.DR.Store(uint32(b))
		return nil
	}
	// No free space in the FIFO. Leave this byte to the ISR.
	waitWriteISR(d, &b, 1)
	return nil
}

func waitWriteISR(d *Driver, p *byte, n int) bool {
	d.wstart = uintptr(unsafe.Pointer(p))
	d.wend = d.wstart + uintptr(n)
	d.wdone.Clear()                    // memory barrier
	internal.AtomicSet(&d.p.IMSC, TXI) // enable Tx FIFO interrupt
	if !d.wdone.Sleep(d.wtimeout) {
		internal.AtomicClear(&d.p.IMSC, TXI)
		return false
	}
	return true
}

// WaitTxDone waits until the last byte (including all the stop bits) from the
// last write operation has been sent. Be aware that your program may not be
// fast enough to perform a time critital action after this function returns.
// Use it only to avoid distrubing your own UART transmision by the following
// actions in program order. For example, disabling the RS485 output driver just
// after returning from WaitTxDone may be done too late because of the other
// gorutines or interrupt handlers so the slave device may start sending a
// response on the bus that is still connected to your RS485 driver (use PIO or
// a hardware solution instaed).
func (d *Driver) WaitTxDone() {
	fr := &d.p.FR
	for fr.LoadBits(BUSY) != 0 {
		runtime.Gosched()
	}
}
