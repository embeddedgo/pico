// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/p/mmap"
)

type Port struct {
	in     mmio.R32[uint32]
	_      [2]uint32
	out    mmio.R32[uint32]
	_      uint32
	outSet mmio.R32[uint32]
	_      uint32
	outClr mmio.R32[uint32]
	_      uint32
	outXor mmio.R32[uint32]
	_      uint32
	oe     mmio.R32[uint32]
	_      uint32
	oeSet  mmio.R32[uint32]
	_      uint32
	oeClr  mmio.R32[uint32]
	_      uint32
	oeXor  mmio.R32[uint32]
}

func P(n int) *Port {
	if uint(n) > 1 {
		return nil
	}
	return (*Port)(unsafe.Pointer(mmap.SIO_BASE + uintptr(n+1)*4))
}

func (p *Port) Bit(n int) Bit {
	if uint(n) > 31 {
		panic("bad GPIO bit number")
	}
	addr := uintptr(unsafe.Pointer(p))
	return Bit{uint8(addr&3<<5 | uintptr(n))}
}
