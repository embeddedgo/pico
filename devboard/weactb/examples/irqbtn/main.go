// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Irqbtn shows how to use interrupts to wait for the onboard button to be
// pressed. See also ../../../pico2/examples/irqpin/main.go.
package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/pico/devboard/weactb/board/buttons"
	"github.com/embeddedgo/pico/devboard/weactb/board/leds"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
)

func main() {
	irq.IO_BANK0.Enable(rtos.IntPrioLow, 0) // enable the IO_BANK0 IRQ on Proc0
	for {
		waitBtn(0)
		waitBtn(1)
		leds.User.Toggle()
	}
}

var note rtos.Note

func waitBtn(state int) {
	pin := buttons.User.Pin()
	for {
		note.Clear()
		pin.SetDstIRQ(iomux.Proc0, iomux.EdgeLow|iomux.EdgeHigh)
		wait := time.Duration(-1)
		if buttons.User.Read() == state {
			wait = 50 * time.Millisecond // we want 50 ms of stable state
		}
		if !note.Sleep(wait) {
			pin.SetDstIRQ(iomux.Proc0, 0)
			return
		}
	}
}

//go:interrupthandler
func IO_BANK0_Handler() {
	pin := buttons.User.Pin()
	pin.SetDstIRQ(iomux.Proc0, 0)
	pin.ClearIRQ(iomux.EdgeLow | iomux.EdgeHigh)
	note.Wakeup()
}
