// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Display draws on the connected display. See also  ../../module/waveshare/examples/pico-lcd-1.3/main.go.
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
	"github.com/embeddedgo/pico/hal/spi/spi1dma"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1

		// This ridiculous pinout makes this example compatible with the
		// Waveshare Pico-LCD-1.3 hat.
		dc   = pins.GP8
		csn  = pins.GP9
		sck  = pins.GP10
		mosi = pins.GP11
		rst  = pins.GP12 // optional, connect to 3V (exception SSD1306)
		//bl   = pins.GP13 // backlight, optional

		// This is the only SPI1 RX capable pin left by Pico-LCD-1.3 (it uses
		// the remaining two for DC and RST).
		miso = pins.GP28_A2
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	// Setup SPI driver
	sm := spi1dma.Master()
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
	//dp := displays.Waveshare_1i3_240x240_IPS_ST7789()

	// Most of the displays accept significant overclocking. The values below
	// work for the ILI9341 controller (80/26.8 MHz calculated, 62.5/20.8 MHz
	// actual values in case of 125 MHz Pico).
	//dp.MaxWriteClk *= 8
	//dp.MaxReadClk *= 4

	dci := tftdci.NewSPI(
		sm, csn, dc,
		spi.CPOL1|spi.CPHA1, // faster than CPOL0|CPHA0 (no gaps between words)
		dp.MaxReadClk, dp.MaxWriteClk,
	)

	fmt.Println("SPI baudrate:")
	fmt.Println("- write:", sm.SetBaudrate(dp.MaxWriteClk))
	fmt.Println("- read: ", sm.SetBaudrate(dp.MaxReadClk))

	fmt.Println("*** Start ***")

	disp := dp.New(dci)
	for {
		examples.RotateDisplay(disp)
		examples.DrawText(disp)
		examples.GraphicsTest(disp)
	}
}
