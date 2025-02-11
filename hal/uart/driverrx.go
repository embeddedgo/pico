// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"runtime"

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
	return e.s[1:]
}

func (d *Driver) Read(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	p := d.p
	for p.FR.LoadBits(RXFE) != 0 {
		runtime.Gosched()
	}
	for p.FR.LoadBits(RXFE) == 0 {
		v := p.DR.Load()
		s[n] = byte(v)
		e := v >> 8 & 15
		n++
		if e != 0 {
			err = &errStr[e-1]
			break
		}
		if n >= len(s) {
			break
		}
	}
	return
}
