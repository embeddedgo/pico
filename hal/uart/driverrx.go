// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
)

func (d *Driver) EnableRx() {
	internal.AtomicSet(&d.p.CR, UARTEN|RXE)
}

func (d *Driver) DisableRx() {
	internal.AtomicClear(&d.p.CR, RXE)
}

type Error struct {
	s string
}

var errStr = [15]Error{
	{"\x01fra"},
	{"\x02par"},
	{"\x03par+fra"},
	{"\x04bre"},
	{"\x05bre+fra"},
	{"\x06bre+par"},
	{"\x07bre+par+fra"},
	{"\x08ove"},
	{"\x09ove+fra"},
	{"\x0Aove+par"},
	{"\x0Bove+par+fra"},
	{"\x0Cove+bre"},
	{"\x0Dove+bre+fra"},
	{"\x0Eove+bre+par"},
	{"\x0Fove+bre+par+fra"},
}

// Status return the error flags as in the RSR register (see FE, PE, BE, OE).
func (e *Error) Status() uint32 {
	return uint32(e.s[0])
}

func (e *Error) Error() string {
	return "uart: " + e.s[1:]
}

func (d *Driver) Read(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	var e uint32
	p := d.p
	if p.FR.LoadBits(RXFE) != 0 {
		// No data in the FIFO. The ISR will read the beggining of new data.
		waitReadISR(d, &s[0], len(s))
		n = len(s) - int(d.rend-d.rstart)
		if e = d.rerr; e != 0 {
			goto error
		}
		if n == len(s) {
			return
		}
		// Check for more data available in the FIFO before return.
	}
	if len(s)-n > 3 && p.RIS.LoadBits(RXI) != 0 {
		// There are at least fifoLen/2 ready bytes in the FIFO.
		fast := s[:min(n+fifoLen/2, len(s))]
		for {
			v := p.DR.Load()
			fast[n] = byte(v)
			n++
			if e = v >> 8 & 15; e != 0 {
				goto error
			}
			if n >= len(fast) {
				break
			}
		}
		if n == len(s) {
			return
		}
	}
	for p.FR.LoadBits(RXFE) == 0 {
		// There is at least 1 ready byte in the FIFO.
		v := p.DR.Load()
		s[n] = byte(v)
		n++
		if e = v >> 8 & 15; e != 0 {
			goto error
		}
		if n >= len(s) {
			break
		}
	}
	return
error:
	err = &errStr[e-1]
	return
}

func (d *Driver) ReadByte() (b byte, err error) {
	var e uint32
	p := d.p
	if p.FR.LoadBits(RXFE) == 0 {
		v := p.DR.Load()
		b = byte(v)
		if e = v >> 8 & 15; e != 0 {
			goto error
		}
		return
	}
	// No data in the FIFO. The ISR will read a byte for us.
	waitReadISR(d, &b, 1)
	if e = d.rerr; e != 0 {
		goto error
	}
	return
error:
	err = &errStr[e-1]
	return
}

func waitReadISR(d *Driver, p *byte, n int) {
	d.rstart = uintptr(unsafe.Pointer(p))
	d.rend = d.rstart + uintptr(n)
	d.rerr = 0
	d.rready.Clear()
	internal.AtomicSet(&d.p.IMSC, RXI|RTI) // enable Rx FIFO interrupts
	d.rready.Sleep(-1)
}
