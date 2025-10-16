// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"embedded/mmio"
	"structs"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/hal/system/clock"
)

type SM struct {
	_ structs.HostLayout

	CLKDIV    mmio.U32
	EXECCTRL  mmio.R32[EXECCTRL]
	SHIFTCTRL mmio.R32[SHIFTCTRL]
	ADDR      mmio.U32
	INSTR     mmio.U32
	PINCTRL   mmio.R32[PINCTRL]
}

func (sm *SM) PIO() *PIO {
	addr := uintptr(unsafe.Pointer(sm)) &^ (pioStep - 1)
	return (*PIO)(unsafe.Pointer(addr))
}

func (sm *SM) Num() int {
	addr := uintptr(unsafe.Pointer(sm))
	return int(addr>>2-2) & (numSM - 1)
}

// Disable disables the state machine (stops executing program).
func (sm *SM) Disable() {
	internal.AtomicClear(&sm.PIO().CTRL, (1<<SM_ENABLEn)<<uint(sm.Num()))
}

// Enable enables the state machine (starts executing program)
func (sm *SM) Enable() {
	internal.AtomicSet(&sm.PIO().CTRL, (1<<SM_ENABLEn)<<uint(sm.Num()))
}

// Reset disables the state machine, resets its internal state and applies the
// default configuration.
func (sm *SM) Reset() {
	pio := sm.PIO()
	smm := CTRL(1) << uint(sm.Num())

	// Disable.
	internal.AtomicClear(&pio.CTRL, smm<<SM_ENABLEn)

	// Default config and clearing FIFOs.
	sm.CLKDIV.Store(1 << INTn)
	sm.EXECCTRL.Store(31 << WRAP_TOPn)
	sm.SHIFTCTRL.Store(FJOIN_RX | OUT_SHIFTDIR | IN_SHIFTDIR)
	sm.SHIFTCTRL.Store(OUT_SHIFTDIR | IN_SHIFTDIR) // join-disjoin clears FIFOs
	sm.PINCTRL.Store(5 << SET_COUNTn)

	// Clear some internal state.
	internal.AtomicSet(&pio.CTRL, smm<<SM_RESTARTn)
	internal.AtomicClear(&pio.CTRL, smm<<SM_RESTARTn)
	internal.AtomicSet(&pio.CTRL, smm<<CLKDIV_RESTARTn)
	internal.AtomicClear(&pio.CTRL, smm<<CLKDIV_RESTARTn)
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
	sm.INSTR.Store(0xe000 + uint32(initPC))
}

// SetClkFreq configures the SM to run at the given frequency. It returns the
// actual frequency which may differ from freq due to rounding. Sea also
// SetClkDiv.
func (sm *SM) SetClkFreq(freq int64) (actual int64) {
	pclk := clock.PERI.Freq()
	div := pclk * 256 / freq
	if div>>24 != 0 {
		return
	}
	sm.CLKDIV.Store(uint32(div) << FRACn)
	return pclk * 256 / div
}

// SetClkDiv configures the SM to run at the clock equal to
// clock.PERI.Freq() * 256 / (divInt * 256 + divFrac). See also SetClkFerq.
func (sm *SM) SetClkDiv(divInt, divFrac uint) {
	sm.CLKDIV.Store(uint32(divInt<<16+divFrac<<8) << FRACn)
}
