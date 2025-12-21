// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system/clock"
)

type SM struct {
	r SMRegs
}

func (sm *SM) PIO() *PIO {
	addr := uintptr(unsafe.Pointer(sm)) &^ (pioStep - 1)
	return (*PIO)(unsafe.Pointer(addr))
}

func (sm *SM) Num() int {
	addr := uintptr(unsafe.Pointer(sm))
	return int(addr>>2-2) & (numSM - 1)
}

func (sm *SM) Regs() *SMRegs {
	return &sm.r
}

// Disable disables the state machine (stops executing program).
func (sm *SM) Disable() {
	internal.AtomicClear(&sm.PIO().p.CTRL, (1<<SM_ENABLEn)<<uint(sm.Num()))
}

// Enable enables the state machine (starts executing program)
func (sm *SM) Enable() {
	internal.AtomicSet(&sm.PIO().p.CTRL, (1<<SM_ENABLEn)<<uint(sm.Num()))
}

// Reset disables the state machine, resets its internal state and applies the
// default configuration.
func (sm *SM) Reset() {
	pio := sm.PIO()
	smm := CTRL(1) << uint(sm.Num())

	// Disable.
	internal.AtomicClear(&pio.p.CTRL, smm<<SM_ENABLEn)

	// Default config and clearing FIFOs.
	sm.r.CLKDIV.Store(1 << INTn)
	sm.r.EXECCTRL.Store(31 << WRAP_TOPn)
	sm.r.SHIFTCTRL.Store(FJOIN_RX | OUT_SHIFTDIR | IN_SHIFTDIR)
	sm.r.SHIFTCTRL.Store(OUT_SHIFTDIR | IN_SHIFTDIR) // join-disjoin clears FIFOs
	sm.r.PINCTRL.Store(5 << SET_COUNTn)

	// Clear some internal state.
	internal.AtomicSet(&pio.p.CTRL, smm<<SM_RESTARTn)
	internal.AtomicClear(&pio.p.CTRL, smm<<SM_RESTARTn)
	internal.AtomicSet(&pio.p.CTRL, smm<<CLKDIV_RESTARTn)
	internal.AtomicClear(&pio.p.CTRL, smm<<CLKDIV_RESTARTn)
}

// Configure configures the state machine to run the program prog starting from
// the instruction in the memory slot initPC. It doesn't reset the state machine
// before applying the program configuration (see SM.Reset). It doesn't load the
// program to the instruction memory (see PIO.Load).
func (sm *SM) Configure(prog Program, initPC int) {
	if uint(initPC) >= imCap {
		panic("pio: bad initPC")
	}
	prog.AlterSM(sm)
	sm.r.INSTR.Store(0xe000 + uint32(initPC))
}

// SetClkFreq configures the SM to run at the given frequency. It returns the
// actual frequency which may differ from the freq due to rounding. Sea also
// SetClkDiv.
func (sm *SM) SetClkFreq(freq int64) (actual int64) {
	pclk := clock.PERI.Freq()
	div := pclk * 256 / freq
	if div>>24 != 0 {
		return
	}
	sm.r.CLKDIV.Store(uint32(div) << FRACn)
	return pclk * 256 / div
}

// SetClkDiv configures the SM to run at the clock equal to
// clock.PERI.Freq() * 256 / (divInt * 256 + divFrac). See also SetClkFerq.
func (sm *SM) SetClkDiv(divInt, divFrac uint) {
	sm.r.CLKDIV.Store(uint32(divInt<<16+divFrac<<8) << FRACn)
}

// SetPinBase sets the base pin for out, set and sideset operations.
func (sm *SM) SetPinBase(out, set, sideset iomux.Pin) {
	gpioBase := int(sm.PIO().p.GPIOBASE.LoadBits(16))
	outBase := PINCTRL(int(out) - gpioBase)
	setBase := PINCTRL(int(set) - gpioBase)
	sidesetBase := PINCTRL(int(sideset) - gpioBase)
	if outBase > 31 || setBase > 31 || sidesetBase > 31 {
		panic("pio: pin out of range")
	}
	sm.r.PINCTRL.StoreBits(
		OUT_BASE|SET_BASE|SIDESET_BASE,
		outBase<<OUT_BASEn|setBase<<SET_BASEn|sidesetBase<<SIDESET_BASEn,
	)
}
