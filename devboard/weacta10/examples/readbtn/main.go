// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/embeddedgo/pico/devboard/weacta10/board/buttons"
	"github.com/embeddedgo/pico/devboard/weacta10/board/leds"
)

func main() {
	for {
		if buttons.User.Read() != 0 {
			leds.Green.SetOn()
			leds.Blue.SetOff()
		} else {
			leds.Blue.SetOn()
			leds.Green.SetOff()
		}
	}
}
