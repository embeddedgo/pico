// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

func main() {
	pin := pins.GP15
	pin.Setup(iomux.Schmitt | iomux.PullUp | iomux.IE)
	gpio.UsePin(pin)
	inp := gpio.BitForPin(pin)

	for {
		if inp.Load() == 0 {
			leds.User.SetOn()
		} else {
			leds.User.SetOff()
		}
	}
}
