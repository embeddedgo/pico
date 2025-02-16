// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Crudeblinky flashes the on-board LED without using HAL. Check ../blinky to
// see a better way to do the same thing.
package main

import (
	"github.com/embeddedgo/pico/p/iobank"
	"github.com/embeddedgo/pico/p/padsbank"
	"github.com/embeddedgo/pico/p/sio"
)

const LEDpin = 25

func main() {
	ledpin := &padsbank.PADS_BANK0().GPIO[LEDpin]
	ledio := &iobank.IO_BANK0().GPIO[LEDpin].CTRL
	sio := sio.SIO()
	oeset := &sio.GPIO_OE_SET
	outset := &sio.GPIO_OUT_SET
	outclr := &sio.GPIO_OUT_CLR

	// The RP2530 pins have isolation latches that prevent changes to the pin
	// state during boot from the deep sleep. The sequence below makes little
	// sense for the onboard LED but demonstrates how to configure the pin to
	// the known state before disabling its isolation latch.

	// Setup SIO. Output the low state on the LED pin..
	outclr.Store(1 << LEDpin)
	oeset.Store(1 << LEDpin)

	// Connect the input and output of the LED pin to the SIO without disabling
	// its isolation latch. Set the output driver strength to 8 mA.
	ledpin.Store(padsbank.ISO | padsbank.IE | padsbank.D8MA)
	ledio.Store(iobank.F5_SIO)

	// Disable the isolation latch.
	ledpin.ClearBits(padsbank.ISO)

	// Blinking in a loop. As the system timer isn't initialized the required
	// delays are implemented by writting GPIO registers multiple times.
	for {
		for i := 0; i < 1e6; i++ {
			outset.Store(1 << LEDpin)
		}
		for i := 0; i < 1e6; i++ {
			outclr.Store(1 << LEDpin)
		}
	}
}
