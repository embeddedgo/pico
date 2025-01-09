// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
func (b Bit) IsValid() bool {
	return b.h>>5 != 0
}

// Port returns the port where the bit is located.
func (b Bit) Port() *Port {
	return (*Port)(unsafe.Pointer(mmap.SIO_BASE + uintptr(b.h)>>5*4))
}

// Num returns the bit number in the port.
func (b Bit) Num() int {
	return int(b.h & 31)
}

// Mask returns a bitmask that represents the bit in a port
func (b Bit) Mask() uint32 {
	return 1 << uint(b.Num())
}

// OutEnabled reports the output enabled state.
func (b Bit) OutEnabled() bool {
	return b.Port().oe.LoadBits(b.Mask()) != 0
}

// EnableOut enables output for this bit.
func (b Bit) EnableOut() {
	b.Port().oeSet.Store(b.Mask())
}

// DisableOut disables output for this bit.
func (b Bit) DisableOut() {
	b.Port().oeClr.Store(b.Mask())
}

// Samples the value of the connected pin.
func (b Bit) Load() int {
	return int(b.Port().in.Load()) >> uint(b.Num()) & 1
}

// LoadOut returns the output value set for this pin.
func (b Bit) LoadOut() int {
	return int(b.Port().out.Load()) >> uint(b.Num()) & 1
}

// Set sets the output value of this bit to 1 in one atomic operation.
func (b Bit) Set() {
	b.Port().outSet.Store(1 << uint(b.Num()))
}

// Clear sets the output value of this bit to 0 in one atomic operation.
func (b Bit) Clear() {
	b.Port().outClr.Store(1 << uint(b.Num()))
}

// Toggle toggles the output value of this bit in one atomic operation.
func (b Bit) Toggle() {
	b.Port().outXor.Store(1 << uint(b.Num()))
}

// Store sets the bit value to the least significant bit of val.
func (b Bit) Store(val int) {
	port := b.Port()
	mask := uint32(1) << uint(b.Num())
	if val&1 != 0 {
		port.outSet.Store(mask)
	} else {
		port.outClr.Store(mask)
	}
}

func UsePin(pin iomux.Pin) {
	pin.SetAltFunc(iomux.GPIO)
}

func BitForPin(pin iomux.Pin) Bit {
	return Bit{uint8(pin + 32)}
}

/*

// Interrupt configuration constants
const (
	IntLow     = 0 // interrupt is low-level sensitive
	IntHigh    = 1 // interrupt is high-level sensitive
	IntRising  = 2 // interrupt is rising-edge sensitive
	IntFalling = 3 // interrupt is falling-edge sensitive
)

// IntConf returns the interrupt configuration of bit.
func (b Bit) IntConf() int {
	n := uint(b.Num())
	shift := n * 2 & 15
	return int(b.Port().IntCfg[n>>4].Load()>>shift) & 3
}

// SetIntConf sets the interrupt configuration for bit.
func (b Bit) SetIntConf(cfg int) {
	n := uint(b.Num())
	shift := n * 2 & 15
	b.Port().IntCfg[n>>4].StoreBits(3<<shift, uint32(cfg<<shift))
}

// IntPending reports whether the interrupt coresponding to b is pending.
func (b Bit) IntPending() bool {
	return b.Port().Pending.LoadBits(1<<uint(b.Num())) != 0
}

// ClearPending clears the pending state of the interrupt coresponding to b.
func (b Bit) ClearPending() {
	b.Port().Pending.Store(1 << uint(b.Num()))
}

// ConnectMux works like Port.ConnectMux(b.Mask())
func (b Bit) ConnectMux() {
	b.Port().ConnectMux(1 << uint(b.Num()))
}

// ConnectMux reports wheter the bit is connected to IOMUX.
func (b Bit) MuxConnected() bool {
	return b.Port().MuxConnected()>>uint(b.Num())&1 != 0
}

*/
