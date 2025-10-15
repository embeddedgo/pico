// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/embeddedgo/display/pix/displays"
	"github.com/embeddedgo/display/pix/examples"

	"github.com/embeddedgo/pico/dci/tftdci"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c0dma"
	"github.com/embeddedgo/pico/hal/system/clock"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		sda   = pins.GP20
		scl   = pins.GP21
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	fmt.Println("Clocks:")
	fmt.Println("REF: ", clock.REF.Freq()/1e6, "MHz")
	fmt.Println("SYS: ", clock.SYS.Freq()/1e6, "MHz")
	fmt.Println("PERI:", clock.PERI.Freq()/1e6, "MHz")
	fmt.Println("HSTX:", clock.HSTX.Freq()/1e6, "MHz")
	fmt.Println("USB: ", clock.USB.Freq()/1e6, "MHz")
	fmt.Println("ADC: ", clock.ADC.Freq()/1e6, "MHz")
	fmt.Println()

	// I2C
	m := i2c0dma.Master()
	m.UsePin(sda, i2c.SDA)
	m.UsePin(scl, i2c.SCL)
	m.Setup(400e3)

	dci := tftdci.NewI2C(m, 0b0111100)
	disp := displays.Adafruit_0i96_128x64_OLED_SSD1306().New(dci)
	for {
		examples.RotateDisplay(disp)
		examples.DrawText(disp)
		examples.GraphicsTest(disp)
	}

}
