// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package segled

import (
	"github.com/embeddedgo/display/segdisp"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
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

// WriteBytes implements segdisp.ShiftReg8 interface.
func (sr *ShiftReg) WriteBytes(p []byte) {
	for _, b := range p {
		for i := 7; i >= 0; i-- {
			bit := int(b) >> uint(i)
			sr.clk.Clear()
			sr.din.Store(bit)
			sr.clk.Set()
		}
	}
}

// Latch implements segdisp.ShiftReg8 interface.
func (sr *ShiftReg) Latch() {
	sr.rclk.Set()
	sr.rclk.Clear()
}

var Display *segdisp.Seg8

func init() {
	tm := segdisp.NewTimeMux8(NewShiftReg(pins.GP11, pins.GP10, pins.GP9), 60)
	tm.Start()
	Display = segdisp.NewSeg8(4, 1, tm)
}
