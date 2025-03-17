// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// WS2812 uses the UART peripheral to drive the string of the WS2812 RGB LEDs.
package main

import (
	"image/color"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/devboard/pico2/board/pwr"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart1"
	"github.com/embeddedgo/rgbled/ws281x/wsuart"
)

func main() {
	// Reduce the noise on the WS2812 data signal.
	//
	// If we use a simple circuit, powering LEDs from VBUS (USB 5V) and
	// connecting the WS2812 data signal directly to the UART TX pin (3.3V),
	// our data signal is already out off spec and any additional ripple can
	// only worse things. If your Pico is powered from USB (VBUS) you can
	// improve things slightly by powering LDEs from VSYS, thanks to the voltage
	// drop on the schottky diode between VBUS and VSYS (5V - 0.3V = 4.7V) but
	// the total current flowing through the diode must be < 1A. Also adding a
	// 200 Ω series resisitor on the data line can reduce oscilations and
	// protect the Pico IO pin a little from 5V in case of failure.
	pwr.SetPowerSave(false) // force the onboard DCDC to work in PWM mode

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
