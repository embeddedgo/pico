// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"runtime"
	"time"

	"github.com/embeddedgo/pico/devboard/weacta10/board/buttons"
	"github.com/embeddedgo/pico/devboard/weacta10/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart1"
	"github.com/embeddedgo/rgbled"
	"github.com/embeddedgo/rgbled/ws281x/wsuart"
)

func main() {
	tx := pins.GP22

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	tx.SetAltFunc(iomux.OutInvert)
	u := uart1.Driver()
	u.UsePin(tx, uart.TXD)
	u.Setup(uart.Word7b, 3_000_000_000/1390)
	u.EnableTx()

	grb := wsuart.GRB
	strip := wsuart.Make(8 * 8)

	for {
		for _, c := range []rgbled.Color{
			rgbled.RGB(127, 0, 0),
			rgbled.RGB(255, 0, 0),
			rgbled.RGB(0, 127, 0),
			rgbled.RGB(0, 255, 0),
			rgbled.RGB(0, 0, 127),
			rgbled.RGB(0, 0, 255),
			rgbled.RGB(127, 127, 0),
			rgbled.RGB(255, 255, 0),
			rgbled.RGB(0, 127, 127),
			rgbled.RGB(0, 255, 255),
			rgbled.RGB(127, 0, 127),
			rgbled.RGB(255, 0, 255),
			rgbled.RGB(127, 127, 127),
			rgbled.RGB(255, 255, 255),
		} {
			pixel := grb.Pixel(c)
			for i := 0; i < 64; i++ {
				strip.Clear()
				strip[i] = pixel
				u.Write(strip.Bytes())
				for buttons.User.Read() == 0 {
					runtime.Gosched()
				}
				time.Sleep(time.Second / 8)
			}
		}
	}
}
