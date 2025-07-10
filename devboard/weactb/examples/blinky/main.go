// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Blinky flashes the both on-board LED.
package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/weactb/board/buttons"
	"github.com/embeddedgo/pico/devboard/weactb/board/leds"
)

func main() {
	for {
		delay := time.Second / 2
		if buttons.User.Read() != 0 {
			delay = time.Second / 6
		}
		leds.User.Toggle()
		time.Sleep(delay)
	}
}
