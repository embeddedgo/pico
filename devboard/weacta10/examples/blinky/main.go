// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Blinky flashes the both on-board LEDs using two goroutines.
package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/devboard/weacta10/board/leds"
)

func blink(led common.LED, period time.Duration) {
	ton := period / 4
	toff := period - ton
	for {
		led.SetOn()
		time.Sleep(ton)
		led.SetOff()
		time.Sleep(toff)
	}
}

func main() {
	go blink(leds.Blue, time.Second)
	blink(leds.Green, time.Second/2)
}
