// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/dma/dmairq"
	"github.com/embeddedgo/pico/hal/i2c"
)

func MasterDMA(p *i2c.Periph) (d *i2c.Master) {
	dma0 := dma.DMA(0)
	ch := dma0.AllocChannel()
	d = i2c.NewMaster(i2c.I2C(0), ch)
	dmairq.SetISR(ch, d.DMAISR)
	return
}
