// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SPI loop test: wire GP3 and GP4 together.
package main

import (
	"fmt"
	"slices"
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
		mosi  = pins.GP3
		miso  = pins.GP4
		csn   = pins.GP5
		sck   = pins.GP6
	)

	// Serial console
	u := uart0.Driver()
	u.UsePin(conTx, uart.TXD)
	u.UsePin(conRx, uart.RXD)
	u.Setup(uart.Word8b, 115200)
	u.EnableTx()
	u.EnableRx()

	// Setup SPI0 driver
	//dma0 := dma.DMA(0)
	sm := spi.NewMaster(spi.SPI(0), dma.Channel{}, dma.Channel{})
	sm.UsePin(miso, spi.RXD)
	sm.UsePin(mosi, spi.TXD)
	sm.UsePin(csn, spi.CSN)
	sm.UsePin(sck, spi.SCK)
	sm.Disable()
	sm.Setup(spi.Word8b, 1e6)

	// Data to sent.
	s8 := ">> 0123456789 abcdefghijklmnoprstuvwxyz ABCDEFGHIJKLMNOPRSTUVWXYZ <<"
	s16 := make([]uint16, 70)
	for i := range s16 {
		s16[i] = uint16(0x9000 + i)
	}
	// Make the receive buffers a little bigger than required to test the
	// returned length.
	buf8 := make([]uint8, len(s8)+3)
	buf16 := make([]uint16, len(s16)+3)

	for {
		sm.SetConfig(spi.Word8b)
		n := sm.WriteStringRead(s8, buf8)
		if s8 == string(buf8[:n]) {
			fmt.Fprint(u, "WriteStringRead ok\r\n")
		} else {
			fmt.Fprintf(u, "WriteStringRead err: '%s'\r\n", buf8[:n])
		}

		for i := range 0x100 {
			b := byte(i)
			if sm.WriteReadByte(b) != b {
				fmt.Fprintf(u, "WriteReadByte err: %x\r\n", b)
				goto ok1
			}
		}
		fmt.Fprint(u, "WriteReadByte ok\r\n")
	ok1:

		sm.SetConfig(spi.Word16b)
		n = sm.WriteRead16(s16, buf16)
		if slices.Equal(s16, buf16[:n]) {
			fmt.Fprint(u, "WriteRead16 ok\r\n")
		} else {
			fmt.Fprintf(u, "WriteRead16 err: %x\r\n", buf16[:n])
		}

		for i := range 0x10000 {
			w := uint16(i)
			if sm.WriteReadWord16(w) != w {
				fmt.Fprintf(u, "WriteReadWord16 err: %x\r\n", w)
				goto ok2
			}
		}
		fmt.Fprint(u, "WriteReadWord16 ok\r\n")
	ok2:

		clear(buf8)
		clear(buf16)
		time.Sleep(time.Second)
	}
}
