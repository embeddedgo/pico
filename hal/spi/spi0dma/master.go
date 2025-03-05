// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi0dma

import (
	_ "unsafe"

	"github.com/embeddedgo/pico/hal/spi"
	"github.com/embeddedgo/pico/hal/spi/internal"
)

var master *spi.Master

// Master returns ready to use driver for SPI master.
func Master() *spi.Master {
	if master == nil {
		master = internal.MasterDMA(spi.SPI(0))
	}
	return master
}
