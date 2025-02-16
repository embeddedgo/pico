// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/system/timer/riscvst"
)

func main() {
	system.SetupPico2_125MHz()
	riscvst.Setup()

	ledPin := iomux.P25
	ledPin.Setup(iomux.D8mA)

	gpio.UsePin(ledPin)
	led := gpio.BitForPin(ledPin)
	led.EnableOut()

	for {
		led.Set()
		time.Sleep(100 * time.Millisecond)
		led.Clear()
		time.Sleep(900 * time.Millisecond)
	}
}
