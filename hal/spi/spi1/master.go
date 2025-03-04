// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi1

import (
	_ "unsafe"

	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/spi"
)

var driver *spi.Master

// Master returns ready to use driver for SPI master.
func Master() *spi.Master {
	if driver == nil {
		driver = spi.NewMaster(spi.SPI(1), dma.Channel{}, dma.Channel{})
		driver.Setup(spi.Word8b, 1e5)
	}
	return driver
}
