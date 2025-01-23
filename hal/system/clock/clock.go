// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the code from the pico-sdk (Copyright (c) 2020 Raspberry
// Pi (Trading) Ltd; SPDX-License-Identifier: BSD-3-Clause) loosely translated
// to Go.

package clock

import (
	_ "unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/clocks"
)

const (
	GPOUT0 = Clock(clocks.GPOUT0)
	GPOUT1 = Clock(clocks.GPOUT1)
	GPOUT2 = Clock(clocks.GPOUT2)
	GPOUT3 = Clock(clocks.GPOUT3)
	REF    = Clock(clocks.REF)
	SYS    = Clock(clocks.SYS)
	PERI   = Clock(clocks.PERI)
	HSTX   = Clock(clocks.HSTX)
	USB    = Clock(clocks.USB)
	ADC    = Clock(clocks.ADC)
	nclk   = ADC + 1
)

type Clock int

func (clk Clock) Freq() (hz int64) {
	if uint(clk) >= uint(nclk) {
		return -1
	}
	return int64(clocksHz[clk])
}

var clocksHz [nclk]uint // we don't expect freqency > 4 GHz

func hasGlitchlessMux(clk int) bool {
	return clk == clocks.SYS || clk == clocks.REF
}

// TODO: export this function
func set(clk int, src, auxsrc clocks.CTRL, freqHz uint, div clocks.DIV) {
	c := &clocks.CLOCKS().CLK[clk]
	// If increasing divisor, set divisor before source.
	if div > c.DIV.Load() {
		c.DIV.Store(div)
	}

	if hasGlitchlessMux(clk) && src == clocks.REF_CLKSRC_CLK_REF_AUX {
		// If switching a glitchless slice (REF or SYS) to an aux source, switch
		// away from aux first to avoid passing glitches when changing aux mux.
		internal.AtomicClear(&c.CTRL, clocks.REF_SRC)
		for c.SELECTED.LoadBits(1) == 0 {
		}
	} else {
		// Disable clock. On REF and SYS this does nothing, all other clocks
		// have the ENABLE bit in the same position.
		internal.AtomicClear(&c.CTRL, clocks.GPOUT_ENABLE)
		if hz := clocksHz[clk]; hz != 0 {
			internal.BusyWaitAtLeastCycles((clocksHz[clocks.SYS]/hz + 1) * 3)
		}
	}

	// Set aux mux first, and then glitchless mux if this clock has one.
	internal.AtomicMod(&c.CTRL, clocks.SYS_AUXSRC, c.CTRL.Load(), auxsrc)
	if hasGlitchlessMux(clk) {
		internal.AtomicMod(&c.CTRL, clocks.REF_SRC, c.CTRL.Load(), src)
		for c.SELECTED.LoadBits(1<<src) == 0 {
		}
	}

	// Enable clock and set divisor.
	internal.AtomicSet(&c.CTRL, clocks.GPOUT_ENABLE)
	c.DIV.Store(div)

	// Store this clock freqency.
	clocksHz[clk] = freqHz
}

//go:linkname set github.com/embeddedgo/pico/hal/system.setClock
