// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// This implementation has smaller storage than the one from bit1.go but inlines
// much worse with go1.22.

package gpio

import (
	"unsafe"

	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/p/mmap"
)

// Bit represents a single bit in a GPIO port.
type Bit struct {
	h uint8
}

// IsValid reports whether b represents a valid bit.
//
//go:nosplit
func (b Bit) IsValid() bool {
	return b.h>>5 != 0
}

// Port returns the port where the bit is located.
//
//go:nosplit
func (b Bit) Port() *Port {
	return (*Port)(unsafe.Pointer(mmap.SIO_BASE + uintptr(b.h)>>5*4))
}

// Num returns the bit number in the port.
//
//go:nosplit
func (b Bit) Num() int {
	return int(b.h & 31)
}

// Mask returns a bitmask that represents the bit in a port
//
//go:nosplit
func (b Bit) Mask() uint32 {
	return 1 << uint(b.Num())
}

// OutEnabled reports the output enabled state.
//
//go:nosplit
func (b Bit) OutEnabled() bool {
	return b.Port().oe.LoadBits(b.Mask()) != 0
}

// EnableOut enables output for this bit.
//
//go:nosplit
func (b Bit) EnableOut() {
	b.Port().oeSet.Store(b.Mask())
}

// DisableOut disables output for this bit.
//
//go:nosplit
func (b Bit) DisableOut() {
	b.Port().oeClr.Store(b.Mask())
}

// Load samples the value of the connected pin.
//
//go:nosplit
func (b Bit) Load() int {
	return int(b.Port().in.Load()) >> uint(b.Num()) & 1
}

// LoadOut returns the output value set for this pin.
//
//go:nosplit
func (b Bit) LoadOut() int {
	return int(b.Port().out.Load()) >> uint(b.Num()) & 1
}

// Set sets the output value of this bit to 1 in one atomic operation.
//
//go:nosplit
func (b Bit) Set() {
	b.Port().outSet.Store(1 << uint(b.Num()))
}

// Clear sets the output value of this bit to 0 in one atomic operation.
//
//go:nosplit
func (b Bit) Clear() {
	b.Port().outClr.Store(1 << uint(b.Num()))
}

// Toggle toggles the output value of this bit in one atomic operation.
//
//go:nosplit
func (b Bit) Toggle() {
	b.Port().outXor.Store(1 << uint(b.Num()))
}

// Store sets the bit value to the least significant bit of val.
//
//go:nosplit
func (b Bit) Store(val int) {
	port := b.Port()
	mask := uint32(1) << uint(b.Num())
	if val&1 != 0 {
		port.outSet.Store(mask)
	} else {
		port.outClr.Store(mask)
	}
}

// Bit returns the n-th bit from the p port.
//
//go:nosplit
func (p *Port) Bit(n int) Bit {
	if uint(n) > 31 {
		panic("bad GPIO bit number")
	}
	addr := uintptr(unsafe.Pointer(p))
	return Bit{uint8(addr&3<<5 | uintptr(n))}
}

// BitForPin returns the GPIO bit that corresponds to the given pin.
//
//go:nosplit
func BitForPin(pin iomux.Pin) Bit {
	return Bit{uint8(pin + 32)}
}
