// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/led"
	"github.com/embeddedgo/led/ws281x/wsuart"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/devboard/pico2/board/pwr"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart1"
)

func main() {
	pwr.SetPowerSave(false)

	tx := pins.GP22
	tx.SetAltFunc(iomux.OutInvert)

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	u := uart1.Driver()
	u.UsePin(tx, uart.TXD)
	u.Setup(uart.Word7b, 3_000_000_000/1390)
	u.EnableTx()

	grb := wsuart.GRB
	strip := wsuart.Make(8 * 8)

	for {
		for _, c := range []led.Color{
			led.RGB(127, 0, 0),
			led.RGB(255, 0, 0),
			led.RGB(0, 127, 0),
			led.RGB(0, 255, 0),
			led.RGB(0, 0, 127),
			led.RGB(0, 0, 255),
			led.RGB(127, 127, 0),
			led.RGB(255, 255, 0),
			led.RGB(0, 127, 127),
			led.RGB(0, 255, 255),
			led.RGB(127, 0, 127),
			led.RGB(255, 0, 255),
			led.RGB(127, 127, 127),
			led.RGB(255, 255, 255),
		} {
			pixel := grb.Pixel(c)
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
}
