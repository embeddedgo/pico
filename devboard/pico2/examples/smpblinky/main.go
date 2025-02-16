// Copyright 20254 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Smpblinky tests the Go and noos schedulers. The onboard LED informs about the
// used CPU: CPU0 turns it off, CPU1 turns it on.
//
// The first part of this test blinks the onboard LED three times, switching the
// running thread between CPUs. The second part spawns two busy threads and
// leaves the CPU selection to the schedulers. Now the onboard LED should bilink
// unevenly if the threads run on both CPUs.
//
// See also weacta10/examples/smpblinky for better version of this program that
// uses two LEDs.
package main

import (
	"embedded/rtos"
	"runtime"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/p/sio"
)

func blinkcpu(period int) {
	CPUID := &sio.SIO().CPUID
	for {
		cpuid := 0
		// Busy wait to make this thread really busy.
		for i := 0; i < period; i++ {
			cpuid += int(CPUID.Load())
		}
		if cpuid < period/2 {
			leds.User.SetOff() // the above loop ran mostly on CPU0
		} else {
			leds.User.SetOn() // the above loop ran mostly on CPU1
		}
		runtime.Gosched()
	}
}

func main() {
	// Test the binding a thread to CPU.
	runtime.LockOSThread()
	CPUID := &sio.SIO().CPUID
	cpu := rtos.ExeCtx(0)
	for i := 0; i < 6; i++ {
		if CPUID.Load() == 0 {
			leds.User.SetOff()
		} else {
			leds.User.SetOn()
		}
		time.Sleep(time.Second)
		cpu ^= 1
		rtos.Bind(cpu)
	}
	rtos.Bind(rtos.NotBound)
	runtime.UnlockOSThread()

	// Check the scheduler with not bound threads.
	go blinkcpu(2e6)
	blinkcpu(3e6)
}
