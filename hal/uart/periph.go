// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import (
	"embedded/mmio"
	"structs"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

type Periph struct {
	_ structs.HostLayout

	DR        mmio.U32
	RSR       mmio.U32
	_         [4]uint32
	FR        mmio.R32[FR]
	_         uint32
	ILPR      mmio.U32
	IBRD      mmio.U32
	FBRD      mmio.U32
	LCR_H     mmio.R32[LCR_H]
	CR        mmio.R32[CR]
	IFLS      mmio.U32
	IMSC      mmio.R32[INT]
	RIS       mmio.R32[INT]
	MIS       mmio.R32[INT]
	ICR       mmio.R32[INT]
	DMACR     mmio.R32[DMACR]
	_         [997]uint32
	PERIPHID0 mmio.U32
	PERIPHID1 mmio.U32
	PERIPHID2 mmio.U32
	PERIPHID3 mmio.U32
	PCELLID0  mmio.U32
	PCELLID1  mmio.U32
	PCELLID2  mmio.U32
	PCELLID3  mmio.U32
}

// UART returns the n-th instance of the UART peripheral.
func UART(n int) *Periph {
	if uint(n) > 1 {
		panic("wrong UART number")
	}
	const base = mmap.UART0_BASE
	const step = mmap.UART1_BASE - mmap.UART0_BASE
	return (*Periph)(unsafe.Pointer(base + uintptr(n)*step))
}

func num(p *Periph) int {
	const step = mmap.UART1_BASE - mmap.UART0_BASE
	return int((uintptr(unsafe.Pointer(p)) - mmap.UART0_BASE) / step)
}

// SetReset allows to assert/deassert the reset signal to the UART peripheral.
func (p *Periph) SetReset(assert bool) {
	internal.SetReset(resets.UART0<<uint(num(p)), assert)
}
