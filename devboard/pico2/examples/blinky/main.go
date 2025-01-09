// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/system/timer/riscvst"
	"github.com/embeddedgo/pico/p/sio"
)

func main() {
	system.SetupPico2_125MHz()
	riscvst.Setup()

	ledpin := iomux.P25
	ledpin.Setup(iomux.D4mA)
	ledpin.SetAltFunc(iomux.F5_SIO)

	sio := sio.SIO()
	sio.GPIO_OE_SET.Store(1 << ledpin)
	outset := &sio.GPIO_OUT_SET
	outclr := &sio.GPIO_OUT_CLR

	for {
		outset.Store(1 << ledpin)
		time.Sleep(100 * time.Millisecond)
		outclr.Store(1 << ledpin)
		time.Sleep(900 * time.Millisecond)
	}
}
