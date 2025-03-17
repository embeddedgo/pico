// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Readpin polls the state of the GP15 pin and shows it using the onboard LED.
package main

import (
	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

func main() {
	pin := pins.GP15
	pin.Setup(iomux.Schmitt | iomux.PullUp | iomux.InpEn)
	inp := gpio.UsePin(pin)

	for {
		if inp.Load() == 0 {
			leds.User.SetOn()
		} else {
			leds.User.SetOff()
		}
	}
}
