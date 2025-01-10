// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import (
	"math/bits"

	"github.com/embeddedgo/pico/hal/iomux"
)

// Bit represents a single bit in a GPIO port.
type Bit struct {
	port *Port
	mask uint32
}

// IsValid reports whether b represents a valid bit.
func (b Bit) IsValid() bool {
	return b.port != nil
}

// Port returns the port where the bit is located.
func (b Bit) Port() *Port {
	return b.port
}

// Num returns the bit number in the port.
func (b Bit) Num() int {
	return bits.TrailingZeros(uint(b.mask))
}

// Mask returns a bitmask that represents the bit in a port
func (b Bit) Mask() uint32 {
	return b.mask
}

// OutEnabled reports the output enabled state.
func (b Bit) OutEnabled() bool {
	return b.port.oe.LoadBits(b.mask) != 0
}

// EnableOut enables output for this bit.
func (b Bit) EnableOut() {
	b.port.oeSet.Store(b.mask)
}

// DisableOut disables output for this bit.
func (b Bit) DisableOut() {
	b.port.oeClr.Store(b.mask)
}

// Load samples the value of the connected pin.
func (b Bit) Load() int {
	return int(b.port.in.Load()) >> uint(b.Num()) & 1
}

// LoadOut returns the output value set for this pin.
func (b Bit) LoadOut() int {
	return int(b.port.out.Load()) >> uint(b.Num()) & 1
}

// Set sets the output value of this bit to 1 in one atomic operation.
func (b Bit) Set() {
	b.port.outSet.Store(b.mask)
}

// Clear sets the output value of this bit to 0 in one atomic operation.
func (b Bit) Clear() {
	b.port.outClr.Store(b.mask)
}

// Toggle toggles the output value of this bit in one atomic operation.
func (b Bit) Toggle() {
	b.port.outXor.Store(b.mask)
}

// Store sets the bit value to the least significant bit of val.
func (b Bit) Store(val int) {
	if val&1 != 0 {
		b.port.outSet.Store(b.mask)
	} else {
		b.port.outClr.Store(b.mask)
	}
}

// Bit returns the n-th bit from the p port.
func (p *Port) Bit(n int) Bit {
	if uint(n) > 31 {
		panic("bad GPIO bit number")
	}
	return Bit{p, uint32(1) << uint(n)}
}

// BitForPin returns the GPIO bit that corresponds to the given pin.
func BitForPin(pin iomux.Pin) Bit {
	return Bit{P(int(pin >> 5)), 1 << uint(pin&31)}
}
