// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"time"

	"github.com/embeddedgo/pico/devboard/cytronmaker/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
	"github.com/embeddedgo/rgbled/ws281x/wsuart"
)

func main() {
	tx := pins.GP28_A2
	tx.SetAltFunc(iomux.OutInvert)

	u := uart0.Driver()
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

	for i := 0; ; i++ {
		pixel := grb.Pixel(colors[i%len(colors)])
		u.Write(pixel.Bytes())
		time.Sleep(time.Second / 4)
	}
}
