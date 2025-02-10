// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"runtime"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
)

// EnableTx enables the Tx part of the UART.
func (d *Driver) EnableTx() {
	internal.AtomicSet(&d.p.CR, UARTEN|TXE)
}

// DisableTx waits for the end of transfer (see Flush) and disables the Tx part
// of the UART.
func (d *Driver) DisableTx() {
	d.WaitTxDone()
	internal.AtomicClear(&d.p.CR, TXE)
}

func (d *Driver) Write(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	p := d.p
	// To avoid interrupts write FIFO in thread mode as much as possible.
	if len(s) > 3 {
		// The overhead required to setup a tight write loop may pay off.
		var m int
		if p.FR.LoadBits(TXFE) != 0 {
			m = fifoLen
		} else if p.RIS.LoadBits(TXI) != 0 {
			m = fifoLen / 2
		}
		if m != 0 {
			// There is at least m free locations in the FIFO.
			n = min(m, len(s))
			for _, b := range s[:n] {
				p.DR.Store(uint32(b))
			}
			if n == len(s) {
				return
			}
		}
	}
	for p.FR.LoadBits(TXFF) == 0 {
		// There is at least 1 free location in the FIFO.
		p.DR.Store(uint32(s[n]))
		if n++; n >= len(s) {
			return
		}
	}
	// The remaining data will be written to the FIFO by the ISR.
	d.wdata = &s[n]
	d.wi = 0
	d.wn = len(s) - n
	d.wdone.Clear()                  // memory barrier
	internal.AtomicSet(&p.IMSC, TXI) // enable Tx FIFO interrupt
	d.wdone.Sleep(-1)
	return len(s), nil
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
	d.wdata = &b
	d.wi = 0
	d.wn = 1
	d.wdone.Clear()                  // memory barrier
	internal.AtomicSet(&p.IMSC, TXI) // enable Tx FIFO interrupt
	d.wdone.Sleep(-1)
	return nil
}

// WaitTxDone waits until the last byte (including all the stop bits) from the
// last write operation has been sent.
func (d *Driver) WaitTxDone() {
	fr := &d.p.FR
	for fr.LoadBits(BUSY) != 0 {
		runtime.Gosched()
	}
}
