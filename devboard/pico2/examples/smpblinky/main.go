// Copyright 20254 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"runtime"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/p/sio"
)

func cpuled() {
	runtime.LockOSThread()
	CPUID := &sio.SIO().CPUID
	for {
		const loopN = 1_000_000
		var cpuid uint32
		for i := 0; i < loopN; i++ {
			cpuid += CPUID.Load()
		}
		if cpuid <= loopN/2 {
			leds.User.SetOff() // the above loop ran mostly on CPU0
		} else {
			leds.User.SetOn() // the above loop ran mostly on CPU1
		}
	}
}

func main() {
	go cpuled()
	cpuled()
}
