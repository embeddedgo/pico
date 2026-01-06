// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"errors"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

type PIO struct {
	p Periph
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

func (pio *PIO) Periph() *Periph {
	return &pio.p
}

// SM returns the n-th state machine of pio.
func (pio *PIO) SM(n int) *SM {
	return (*SM)(unsafe.Pointer(&pio.p.SM[n]))
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
	if end > len(pio.p.INSTR_MEM) {
		return 0, errors.New("pio: out off instruction memory")
	}
	im := pio.p.INSTR_MEM[pos:end]
	prog.LoadTo(im)
	gpioBase := pio.p.GPIOBASE.LoadBits(16)
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
