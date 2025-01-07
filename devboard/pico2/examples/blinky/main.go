// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/system/timer/riscvst"
	"github.com/embeddedgo/pico/p/iobank"
	"github.com/embeddedgo/pico/p/padsbank"
	"github.com/embeddedgo/pico/p/sio"
)

const LEDpin = 25

func main() {
	system.SetupPico2_150MHz()
	riscvst.Setup()

	ledpin := &padsbank.PADS_BANK0().GPIO[LEDpin]
	ledio := &iobank.IO_BANK0().GPIO[LEDpin].CTRL
	sio := sio.SIO()
	oeset := &sio.GPIO_OE_SET
	outset := &sio.GPIO_OUT_SET
	outclr := &sio.GPIO_OUT_CLR

	ledpin.Store(padsbank.IE | padsbank.D4MA)
	ledio.Store(iobank.F5_SIO)
	oeset.Store(1 << LEDpin)

	for {
		outset.Store(1 << LEDpin)
		time.Sleep(time.Second / 16)
		outclr.Store(1 << LEDpin)
		time.Sleep(time.Second)
	}
}
