// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Blinky flashes the on-board LED.
package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/core2350b/board/leds"
)

func main() {
	for {
		leds.User.Toggle()
		time.Sleep(time.Second / 2)
	}
}
