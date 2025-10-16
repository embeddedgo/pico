// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"embedded/mmio"
	"structs"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

type PIO struct {
	_ structs.HostLayout

	CTRL              mmio.R32[CTRL]
	FSTAT             mmio.R32[FSTAT]
	FDEBUG            mmio.R32[FDEBUG]
	FLEVEL            mmio.R32[FLEVEL]
	TXF               [4]mmio.U32
	RXF               [4]mmio.U32
	_                 [2]uint32
	INPUT_SYNC_BYPASS mmio.U32
	DBG_PADOUT        mmio.U32
	DBG_PADOE         mmio.U32
	DBG_CFGINFO       mmio.R32[DBG_CFGINFO]
	INSTR_MEM         [32]mmio.U32
	SM                [4]SM
	RXF_PUTGET        [4][4]mmio.U32
	_                 [12]uint32
	GPIOBASE          mmio.U32
	INTR              mmio.R32[INTR]
	IRQ               [2]SIRQ
}

type SIRQ struct {
	_ structs.HostLayout

	E mmio.R32[INTR]
	F mmio.R32[INTR]
	S mmio.R32[INTR]
}

const pioStep = 0x100000

// Block returns the n-th instance of the Programable IO Block
func Block(n int) *PIO {
	if uint(n) > 2 {
		panic("wrong PIO number")
	}
	addr := mmap.PIO0_BASE + uintptr(n)*pioStep
	return (*PIO)(unsafe.Pointer(addr))
}

func (pio *PIO) Num() int {
	addr := uintptr(unsafe.Pointer(pio))
	return int((addr - mmap.PIO0_BASE) / pioStep)
}

func (pio *PIO) SetReset(assert bool) {
	internal.SetReset(resets.PIO0<<uint(pio.Num()), assert)
}
