// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
)

const irqPin = pins.GP15

func main() {
	irqPin.Setup(iomux.Schmitt | iomux.PullUp | iomux.InpEn)
	irqPin.SetDstIRQ(iomux.Proc0, iomux.EdgeLow)
	irq.IO_BANK0.Enable(rtos.IntPrioLow, 0)

	for {
		time.Sleep(2 * time.Second)
		leds.User.Toggle()
		time.Sleep(50 * time.Millisecond)
		leds.User.Toggle()
	}
}

//go:interrupthandler
func IO_BANK0_Handler() {
	irqPin.ClearIRQ(iomux.EdgeLow)
	leds.User.Toggle()
}
