// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package display

import (
	"time"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/displays"
	"github.com/embeddedgo/pico/dci/tftdci"
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/spi"
	"github.com/embeddedgo/pico/hal/spi/spi1dma"
	_ "github.com/embeddedgo/pico/hal/system/init"
)

const Backlight = common.LED(iomux.P13) // TODO: PWM

var LCD *pix.Display

func init() {
	common.ConnectLED(Backlight.Pin(), iomux.D4mA, 0)
	Backlight.SetOn()

	dc = iomux.P08
	csn = iomux.P09
	sck = iomux.P10
	mosi = iomux.P11
	rst = iomux.P12

	// Setup SPI driver
	sm := spi1dma.Master()
	sm.UsePin(miso, spi.RXD)
	sm.UsePin(mosi, spi.TXD)
	sm.UsePin(sck, spi.SCK)

	// Reset the display controller
	reset := gpio.UsePin(rst)
	reset.EnableOut()
	reset.Clear()         // set reset initial steate low
	rst.Setup(iomux.D4mA) // set the rst pin as output
	time.Sleep(time.Millisecond)
	reset.Set()

	dp := displays.Waveshare_1i3_240x240_IPS_ST7789()
	dci := tftdci.NewSPI(
		sm, csn, dc,
		spi.CPOL1|spi.CPHA1, // faster than CPOL0,CPHA0 (no gaps between words)
		dp.MaxReadClk, dp.MaxWriteClk,
	)
	LCD = dp.New(dci)
}
