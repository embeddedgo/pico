// Copyright 2026 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oled

import (
	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/display/pix/displays"
	"github.com/embeddedgo/pico/dci/tftdci"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c1dma"
	"github.com/embeddedgo/pico/hal/iomux"
)

var Display *pix.Display

func init() {
	// Setup I2C
	m := i2c1dma.Master()
	m.UsePin(iomux.P18, i2c.SDA)
	m.UsePin(iomux.P19, i2c.SCL)
	m.Setup(400e3)

	dci := tftdci.NewI2C(m, 0b0111100)
	Display = displays.GNPE_0i41_72x40_OLED_SSD1315.New(dci)
}
