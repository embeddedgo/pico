// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Display draws on the connected display.
package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/display/pix/displays"
	"github.com/embeddedgo/display/pix/examples"

	"github.com/embeddedgo/pico/dci/tftdci"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/spi"
	"github.com/embeddedgo/pico/hal/spi/spi0"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		mosi  = pins.GP3
		miso  = pins.GP4
		csn   = pins.GP5
		sck   = pins.GP6
		dc    = pins.GP7
		rst   = pins.GP8 // optional, connect to 3V (exception SSD1306)
	)

	// Serial console
	u := uart0.Driver()
	u.UsePin(conTx, uart.TXD)
	u.UsePin(conRx, uart.RXD)
	u.Setup(uart.Word8b, 115200)
	u.EnableTx()
	u.EnableRx()

	// Setup SPI0 driver
	sm := spi0.Master()
	sm.UsePin(miso, spi.RXD)
	sm.UsePin(mosi, spi.TXD)
	sm.UsePin(sck, spi.SCK)

	// Reset the display controller (optional, exception SSD1306).
	reset := gpio.UsePin(rst)
	reset.EnableOut()
	reset.Clear()         // set reset initial steate low
	rst.Setup(iomux.D4mA) // set the rst pin as output
	time.Sleep(time.Millisecond)
	reset.Set()

	//dp := displays.Adafruit_0i96_128x64_OLED_SSD1306()
	//dp := displays.Adafruit_1i5_128x128_OLED_SSD1351()
	//dp := displays.Adafruit_1i54_240x240_IPS_ST7789()
	dp := displays.Adafruit_2i8_240x320_TFT_ILI9341()
	//dp := displays.ERTFTM_1i54_240x240_IPS_ST7789()
	//dp := displays.MSP4022_4i0_320x480_TFT_ILI9486()
	//dp := displays.Waveshare_1i5_128x128_OLED_SSD1351()

	// Most of the displays accept significant overclocking.
	//dp.MaxReadClk *= 2
	//dp.MaxWriteClk *= 2

	dci := tftdci.NewSPI(
		spi0.Master(),
		csn, dc,
		spi.CPOL0|spi.CPHA0,
		dp.MaxReadClk, dp.MaxWriteClk,
	)

	fmt.Fprintln(u, "*** Start ***")

	disp := dp.New(dci)
	for {
		examples.RotateDisplay(disp)
		examples.DrawText(disp)
		examples.GraphicsTest(disp)
	}
}
