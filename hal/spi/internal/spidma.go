// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/dma/dmairq"
	"github.com/embeddedgo/pico/hal/spi"
)

func MasterDMA(p *spi.Periph) (d *spi.Master) {
	dma0 := dma.DMA(0)
	rdma := dma0.AllocChannel()
	wdma := dma0.AllocChannel()
	d = spi.NewMaster(p, rdma, wdma)
	d.Setup(spi.Word8b, 1e5) // make this driver somewhat ready to use
	dmairq.SetISR(rdma, d.DMAISR)
	dmairq.SetISR(wdma, d.DMAISR)
	return
}
