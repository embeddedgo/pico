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
	d.Flush()
	internal.AtomicClear(&d.p.CR, TXE)
}

func (d *Driver) Write(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	p := d.p
	// To avoid interrupts first try to write directly into the FIFO.
	fr := p.FR.Load()
	if fr&TXFE != 0 {
		// Empty FIFO
		n = min(fifoLen, len(s))
		for _, b := range s[:n] {
			p.DR.Store(uint32(b))
		}
	} else if fr&TXFF == 0 {
		for n < len(s) {
			p.DR.Store(uint32(s[n]))
			n++
			if p.FR.LoadBits(TXFF) != 0 {
				break
			}
		}
	}
	if n == len(s) {
		return
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
	d.wbyte = b
	d.wdata = &d.wbyte
	d.wi = 0
	d.wn = 1
	d.wdone.Clear()                  // memory barrier
	internal.AtomicSet(&p.IMSC, TXI) // enable Tx FIFO interrupt
	d.wdone.Sleep(-1)
	return nil
}

// Flush waits until the last byte (including all the stop bits) from the last
// write operation has been sent.
func (d *Driver) Flush() {
	fr := &d.p.FR
	for fr.LoadBits(BUSY) != 0 {
		runtime.Gosched()
	}
}
