// Copyright 20254 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Smpblinky tests the Go and noos schedulers. The onboard LEDs inform about
// the used CPU.
//
// There are two goroutines. Each one has its own LED and sets it on if it
// mostly run on CPU1 and off otherwise. Even if running constantly on the same
// CPU the goroutine also blinks its LED shortly to make a sign of life.
package main

import (
	"runtime"

	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/devboard/weacta10/board/leds"
	"github.com/embeddedgo/pico/p/sio"
)

func blinkcpu(period int, led common.LED) {
	CPUID := &sio.SIO().CPUID
	for {
		delay := period / 8
		if led.Get() == 0 {
			delay /= 2 // blinking on is much more visible than blinking off
		}
		cpuid := 0
		// Busy wait to make this thread really busy.
		for i := 0; i < period; i++ {
			cpuid += int(CPUID.Load())

			// Blink shortly to make a sign of life.
			if i == 0 || i == delay {
				led.Toggle()
			}
		}
		if cpuid < period/2 {
			led.SetOff() // the above loop ran mostly on CPU0
		} else {
			led.SetOn() // the above loop ran mostly on CPU1
		}
		runtime.Gosched()
	}
}

func main() {
	go blinkcpu(4.1e6, leds.Blue)
	blinkcpu(4e6, leds.Green)
}
