// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c0

import (
	_ "unsafe"

	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/i2c"
)

var driver *spi.Master

// Master returns ready to use driver for SPI master.
func Master() *spi.Master {
	if driver == nil {
		driver = spi.NewMaster(spi.SPI(0), dma.Channel{}, dma.Channel{})
		driver.Setup(spi.Word8b, 1e5) // make this driver somewhat ready to use
	}
	return driver
}
