// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SPI loop test: wire GP3 and GP4 together.
package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/spi"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
)

var run bool

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		miso  = pins.GP16
		csn   = pins.GP17
		sck   = pins.GP18
		mosi  = pins.GP19
	)

	// Serial console
	u := uart0.Driver()
	u.UsePin(conTx, uart.TXD)
	u.UsePin(conRx, uart.RXD)
	u.Setup(uart.Word8b, 115200)
	u.EnableTx()
	u.EnableRx()

	// Setup SPI0 driver
	dma0 := dma.DMA(0)
	sm := spi.NewMaster(spi.SPI(0), dma0.AllocChannel(), dma0.AllocChannel())
	sm.UsePin(miso, spi.RXD)
	sm.UsePin(mosi, spi.TXD)
	sm.UsePin(csn, spi.CSN)
	sm.UsePin(sck, spi.SCK)
	sm.Disable()
	spiBaud := sm.Setup(spi.Word8b, 1e6)
	p := sm.Periph()
	for tx := uint32(0); ; tx = (tx + 1) & 0xff {
		for p.SR.LoadBits(spi.TNF) == 0 {
		}
		p.DR.Store(tx)
		for p.SR.LoadBits(spi.RNE) == 0 {
		}
		rx := p.DR.Load()
		fmt.Fprintf(
			u, "spiBaud=%d tx=%08b rx=%08b\r\n",
			spiBaud, tx, rx,
		)
		time.Sleep(time.Second)
	}
}
