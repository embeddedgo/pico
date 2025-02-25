// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

type Periph struct {
	CR0       mmio.R32[CR0]
	CR1       mmio.R32[CR1]
	DR        mmio.U32
	SR        mmio.R32[SR]
	CPSR      mmio.U32
	IMSC      mmio.R32[INT]
	RIS       mmio.R32[INT]
	MIS       mmio.R32[INT]
	ICR       mmio.R32[INT]
	DMACR     mmio.R32[DMACR]
	_         [1006]uint32
	PERIPHID0 mmio.U32
	PERIPHID1 mmio.U32
	PERIPHID2 mmio.U32
	PERIPHID3 mmio.U32
	PCELLID0  mmio.U32
	PCELLID1  mmio.U32
	PCELLID2  mmio.U32
	PCELLID3  mmio.U32
}

// SPI returns the n-th SPI peripheral.
func SPI(n int) *Periph {
	if uint(n) > 1 {
		panic("wrong SPI number")
	}
	const base = mmap.SPI0_BASE
	const step = mmap.SPI1_BASE - mmap.SPI0_BASE
	return (*Periph)(unsafe.Pointer(base + uintptr(n)*step))
}

func num(p *Periph) int {
	const step = mmap.SPI1_BASE - mmap.SPI0_BASE
	return int((uintptr(unsafe.Pointer(p)) - mmap.SPI0_BASE) / step)
}

// SetReset allows to assert/deassert the reset signal to the SPI peripheral.
func (p *Periph) SetReset(assert bool) {
	internal.SetReset(resets.SPI0<<uint(num(p)), assert)
}
