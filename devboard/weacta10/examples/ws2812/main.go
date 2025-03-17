// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"time"

	"github.com/embeddedgo/pico/devboard/weacta10/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart1"
	"github.com/embeddedgo/rgbled/ws281x/wsuart"
)

// WS2812 uses the UART peripheral to drive the string of the WS2812 RGB LEDs.
func main() {
	tx := pins.GP22

	tx.SetAltFunc(iomux.OutInvert)
	u := uart1.Driver()
	u.UsePin(tx, uart.TXD)
	u.Setup(uart.Word7b, wsuart.BaudWS2812)
	u.EnableTx()

	colors := []color.RGBA{
		{127, 0, 0, 255},
		{255, 0, 0, 255},
		{0, 127, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 127, 255},
		{0, 0, 255, 255},
		{127, 127, 0, 255},
		{255, 255, 0, 255},
		{0, 127, 127, 255},
		{0, 255, 255, 255},
		{127, 0, 127, 255},
		{255, 0, 255, 255},
		{127, 127, 127, 255},
		{255, 255, 255, 255},
	}
	grb := wsuart.GRB
	strip := wsuart.Make(8 * 8)

	for i := 0; ; i++ {
		pixel := grb.Pixel(colors[i%len(colors)])
		for i := 0; i < 64; i += 8 {
			strip.Clear()
			for k := i; k < i+8; k++ {
				strip[k] = pixel
			}
			u.Write(strip.Bytes())
			time.Sleep(time.Second / 2)
		}
	}
}
