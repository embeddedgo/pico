// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/pico/devboard/cytronmaker/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
	"github.com/embeddedgo/rgbled"
	"github.com/embeddedgo/rgbled/ws281x/wsuart"
)

func main() {
	tx := pins.GP28_A2

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	tx.SetAltFunc(iomux.OutInvert)

	u := uart0.Driver()
	u.UsePin(tx, uart.TXD)
	u.Setup(uart.Word7b, 3_000_000_000/1390)
	u.EnableTx()
	irq.UART0.Enable(rtos.IntPrioLow, 0)

	colors := []rgbled.Color{
		rgbled.RGB(255, 0, 0),
		rgbled.RGB(0, 255, 0),
		rgbled.RGB(0, 0, 255),
		rgbled.RGB(255, 255, 0),
		rgbled.RGB(0, 255, 255),
		rgbled.RGB(255, 0, 255),
		rgbled.RGB(255, 255, 255),
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
	}
}
