// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gpiotest tests all GPIOs exposed as the Pico 2 GPIO pins. It is intended to
// run with the expansion board or test module that provides LEDs on all
// available pins (e.g. Cytron Maker).
package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

func main() {
	// Configure all available pins as GPIO.
	for pin := pins.GP0; pin <= pins.GP22; pin++ {
		pin.Setup(iomux.D4mA)
		pin.SetAltFunc(iomux.GPIO)
	}
	for pin := leds.User.Pin(); pin <= pins.GP28_A2; pin++ {
		pin.Setup(iomux.D4mA)
		pin.SetAltFunc(iomux.GPIO)
	}
	// Available pins + LED as GPIO bitmask
	const pins uint32 = 0b0001_1110_0111_1111_1111_1111_1111_1111

	p0 := gpio.P(0)    // GPIO port 0 controls the pins/pads from 0 to 31
	p0.EnableOut(pins) // enable GPIO output on available pins
	p0.Clear(pins)     // set all pins to the low state

	// Blink all pin LEDs on the expatnsion board.
	for {
		for pin := uint32(1); pin != 0; pin <<= 1 {
			if pin&pins == 0 {
				continue
			}
			for range 4 {
				p0.Toggle(pin)
				time.Sleep(time.Second / 4)
			}
		}
		for range 8 {
			p0.Toggle(pins)
			time.Sleep(time.Second / 4)
		}
	}
}
