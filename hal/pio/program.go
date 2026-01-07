// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import "embedded/mmio"

// A Program is an interface to the loadable PIO program independent of its
// underlying implementation or encoding.
type Program interface {
	// Origin returns the load address required by the program or -1 if not
	// specified.
	Origin() int

	// Len returns the number of the instruction memory slots the program takes
	// when loaded.
	Len() int

	// LoadTo loads the program instructions to the provided instruction memory.
	// The program may assume that the im has required length.
	LoadTo(im []mmio.R32[uint32])

	// AlterSM changes the state machine configuration according to the program
	// directives and/or other settings specific to the program implementation.
	// Generally it doesn't reset the state machine before aplaying changes so
	// some implementation specific part of the current configuration/state will
	// stay unaffected.
	AlterSM(sm *SM)
}

// A StringProgram represents an immutable PIO program stored in a string.
//
// The format is:
//
//	 0: load address (origin),
//	 1: CLKDIV, bytes 1-3
//	 4: EXECTRL, bytes 0-3
//	 8: SHIFTCTRL, bytes 0-3
//	12: PINCTRL, bytes 2-3
//	14: first instruction, bytes 0-1
//	16: second instruction, bytes 0-1
//	...
//
// The multi-byte numbers are stored with the least significant byte first.
//
// The program length is an even number greather than 14. If we need a new
// format in the future while maintaining support for the previous one the new
// format will have an odd length.
type StringProgram string

func (p StringProgram) Origin() int {
	if len(p) <= 14 || len(p)&1 != 0 {
		return -1
	}
	return int(int8(p[0]))
}

func (p StringProgram) Len() int {
	if len(p) <= 14 || len(p)&1 != 0 {
		return 0 // something is wrong with our program
	}
	return (len(p) - 14) >> 1
}

func (p StringProgram) LoadTo(im []mmio.R32[uint32]) {
	n := p.Len()
	for i := range n {
		l := p[14+2*i]
		h := p[15+2*i]
		im[i].Store(uint32(h)<<8 | uint32(l))
	}
}

func (p StringProgram) AlterSM(sm *SM) {
	cd := uint32(p[1])<<8 | uint32(p[2])<<16 | uint32(p[3])<<24
	ec := EXECCTRL(p[4]) | EXECCTRL(p[5])<<8 | EXECCTRL(p[6])<<16 | EXECCTRL(p[7])<<24
	sc := SHIFTCTRL(p[8]) | SHIFTCTRL(p[9])<<8 | SHIFTCTRL(p[10])<<16 | SHIFTCTRL(p[11])<<24
	pc := PINCTRL(p[12])<<16 | PINCTRL(p[13])<<24
	sr := sm.Regs()
	if cd != 1<<16 {
		// Pico-SDK doesn't generate configuration code for CLKDIV=1.0.
		sr.CLKDIV.Store(cd)
	}
	ecm := EXECCTRL(0xffff_ffff)
	if ec&STATUS_SEL == STATUS_SEL {
		ecm &^= STATUS_SEL | STATUS_N
	}
	scm := SHIFTCTRL(0xffff_c01f)
	if icm := IN_COUNT<<1 | 1<<IN_COUNTn; sc&icm == icm {
		scm &^= icm
	}
	pcm := PINCTRL(0xfff0_0000)
	if pc&SIDESET_COUNT == SIDESET_COUNT {
		pcm &^= SIDESET_COUNT
	}
	if pc&SET_COUNT == SET_COUNT {
		pcm &^= SET_COUNT
	}
	if pc&OUT_COUNT == OUT_COUNT {
		pcm &^= OUT_COUNT
	}
	sr.EXECCTRL.StoreBits(ecm, ec)
	sr.SHIFTCTRL.StoreBits(scm, sc)
	sr.PINCTRL.StoreBits(pcm, pc)
}
