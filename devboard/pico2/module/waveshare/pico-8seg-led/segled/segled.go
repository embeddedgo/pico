// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package segled

import (
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

type ShiftReg struct {
	din, clk, rclk gpio.Bit
}

func NewShiftReg(din, clk, rclk iomux.Pin) *ShiftReg {
	din.Setup(iomux.D4mA)
	clk.Setup(iomux.D4mA)
	rclk.Setup(iomux.D4mA)
	gpio.UsePin(din)
	gpio.UsePin(clk)
	gpio.UsePin(rclk)
	sr := &ShiftReg{
		din:  gpio.BitForPin(din),
		clk:  gpio.BitForPin(clk),
		rclk: gpio.BitForPin(rclk),
	}
	sr.din.EnableOut()
	sr.clk.EnableOut()
	sr.rclk.EnableOut()
	return sr
}

func (sr *ShiftReg) WriteByte(b byte) {
	for i := 7; i >= 0; i-- {
		v := int(b) >> uint(i)
		sr.clk.Clear()
		sr.din.Store(v)
		sr.clk.Set()
	}
}

func (sr *ShiftReg) Latch() {
	sr.rclk.Set()
	sr.rclk.Clear()
}
