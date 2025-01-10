// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/p/mmap"
)

// Port represents an 32-bit GPIO port.
type Port struct {
	// go1.22 inilining works much better for mmio.U32 than for mmio.R32[uint32]
	in     mmio.U32
	_      [2]uint32
	out    mmio.U32
	_      uint32
	outSet mmio.U32
	_      uint32
	outClr mmio.U32
	_      uint32
	outXor mmio.U32
	_      uint32
	oe     mmio.U32
	_      uint32
	oeSet  mmio.U32
	_      uint32
	oeClr  mmio.U32
	_      uint32
	oeXor  mmio.U32
}

// P returns the n-th GPIO port.
func P(n int) *Port {
	if uint(n) > 1 {
		return nil
	}
	return (*Port)(unsafe.Pointer(mmap.SIO_BASE + uintptr(n+1)*4))
}

// Num returns the port number.
func (p *Port) Num() int {
	a := uintptr(unsafe.Pointer(p))
	return int((a-mmap.SIO_BASE)/4 - 1)
}

// OutEnabled returns the bits in the output enable state.
func (p *Port) OutEnabled() uint32 {
	return p.oe.Load()
}

// EnableOut enables output for given bits.
func (p *Port) EnableOut(bits uint32) {
	p.oeSet.Store(bits)
}

// DisableOut disables output for given bits.
func (p *Port) DisableOut(bits uint32) {
	p.oeClr.Store(bits)
}

// Load samples the value of the pins connected to this port.
func (p *Port) Load() uint32 {
	return p.in.Load()
}

// LoadOut returns the output value for this port.
func (p *Port) LoadOut() uint32 {
	return p.out.Load()
}

// Set sets the output value for the given bits of this port to 1.
func (p *Port) Set(bits uint32) {
	p.outSet.Store(bits)
}

// Clear sets the output value for the given bits of this port to 0.
func (p *Port) Clear(bits uint32) {
	p.outClr.Store(bits)
}

// Toggle toggles the output value for the given bits of this port.
func (p *Port) Toggle(bits uint32) {
	p.outXor.Store(bits)
}

// Store sets the output of this port to bits.
func (p *Port) Store(bits uint32) {
	p.out.Store(bits)
}

// UsePin connects pin to the GPIO (SIO) peripheral.
func UsePin(pin iomux.Pin) {
	pin.SetAltFunc(iomux.GPIO)
}
