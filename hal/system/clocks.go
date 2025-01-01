// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the code from the pico-sdk (Copyright (c) 2020 Raspberry
// Pi (Trading) Ltd; SPDX-License-Identifier: BSD-3-Clause) loosely translated
// to Go.

// TODO: this probably should be in the hal/clocks package to allow check clock
// frequency and enable/disable some clocks - privileged level requireed?

package system

import (
	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/clocks"
)

var clocksHz [clocks.ADC + 1]uint // we don't expect clock freqency > 4 GHz

func hasGlitchlessMux(clk int) bool {
	return clk == clocks.SYS || clk == clocks.REF
}

func setClock(clk int, src, auxsrc clocks.CTRL, freqHz uint, div clocks.DIV) {
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
	internal.AtomicXor(&c.CTRL, (c.CTRL.Load()^auxsrc)&clocks.SYS_AUXSRC)
	if hasGlitchlessMux(clk) {
		internal.AtomicXor(&c.CTRL, (c.CTRL.Load()^src)&clocks.REF_SRC)
		for c.SELECTED.LoadBits(1<<src) == 0 {
		}
	}

	// Enable clock and set divisor.
	internal.AtomicSet(&c.CTRL, clocks.GPOUT_ENABLE)
	c.DIV.Store(div)

	// Store this clock freqency.
	clocksHz[clk] = freqHz
}
