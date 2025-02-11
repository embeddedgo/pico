// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/led"
	"github.com/embeddedgo/led/ws281x/wsuart"
	"github.com/embeddedgo/pico/devboard/cytronmaker/board/leds"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/uart"
)

var u *uart.Driver

func main() {
	tx := iomux.P28

	tx.Setup(iomux.D2mA)
	tx.SetAltFunc(iomux.UART | iomux.OutInvert)

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	u = uart.NewDriver(uart.UART(0))
	u.Setup(uart.Word7b, 3000000000/1390)
	u.EnableTx()
	irq.UART0.Enable(rtos.IntPrioLow, 0)

	colors := []led.Color{
		led.RGB(255, 0, 0),
		led.RGB(0, 255, 0),
		led.RGB(0, 0, 255),
		led.RGB(255, 255, 0),
		led.RGB(0, 255, 255),
		led.RGB(255, 0, 255),
		led.RGB(255, 255, 255),
	}

	grb := wsuart.GRB
	c1 := colors[0]
	for i := 1; ; i++ {
		c2 := colors[i%len(colors)]
		for a := 0; a < 256; a++ {
			pixel := grb.Pixel(c1.Blend(c2, uint8(a)))
			u.Write(pixel.Bytes())
			time.Sleep(time.Second / 30)
		}
		c1 = c2
		leds.User.Toggle()
	}
}

//go:interrupthandler
func UART0_Handler() {
	u.ISR()
}
