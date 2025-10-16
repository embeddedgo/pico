// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"embedded/mmio"
	"errors"
	"structs"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

const (
	imCap = 32
	numSM = 4
)

type PIO struct {
	_ structs.HostLayout

	CTRL              mmio.R32[CTRL]
	FSTAT             mmio.R32[FSTAT]
	FDEBUG            mmio.R32[FDEBUG]
	FLEVEL            mmio.R32[FLEVEL]
	TXF               [numSM]mmio.U32
	RXF               [numSM]mmio.U32
	_                 [2]uint32
	INPUT_SYNC_BYPASS mmio.U32
	DBG_PADOUT        mmio.U32
	DBG_PADOE         mmio.U32
	DBG_CFGINFO       mmio.R32[DBG_CFGINFO]
	INSTR_MEM         [imCap]mmio.R32[uint32]
	SM                [numSM]SM
	RXF_PUTGET        [numSM][4]mmio.U32
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

// Load loads the PIO program into instruction memory starting at the given
// position pos. If the program requires a specific location in the instruction
// memory (encoded in the program itself) the pos must be -1.
func (pio *PIO) Load(prog Program, pos int) (actualPos int, err error) {
	origin := prog.Origin()
	if pos == -1 {
		pos = origin
	} else if origin != -1 {
		return 0, errors.New("pio: non-relocatable program")
	}
	if pos == -1 {
		pos = 0 // TODO: find a free chunk of the instruction memory
	}
	end := pos + prog.Len()
	if end > len(pio.INSTR_MEM) {
		return 0, errors.New("pio: out off instruction memory")
	}
	im := pio.INSTR_MEM[pos:end]
	prog.LoadTo(im)
	gpioBase := pio.GPIOBASE.LoadBits(16)
	for i := range im {
		op := im[i].Load()
		switch op & 0xe000 {
		case 0x2000: // wait
			if gpioBase != 0 && op&(3<<5) == 0 {
				// Fix the pin number in the wait gpio instruction.
				im[i].Store(op ^ gpioBase)
			}
		case 0xe000: // jmp
			// Update the jump adresses according to the program position.
			im[i].Store(op + uint32(pos))
		}
	}
	return pos, nil
}
