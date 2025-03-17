// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Blinky flashes the on-board LED. It also shows how to configure and connect
// the pin to the GPIO peripheral. Check ../blinky to see a simpler version of
// this program that uses the board/leds package.
package main

import (
	"time"

	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

func main() {
	ledPin := iomux.P25
	ledPin.Setup(iomux.D8mA)

	led := gpio.UsePin(ledPin)
	led.EnableOut()

	for {
		led.Set()
		time.Sleep(100 * time.Millisecond)
		led.Clear()
		time.Sleep(900 * time.Millisecond)
	}
}
