// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SPI loop test: wire GP3 and GP4 together.
package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/embeddedgo/pico/hal/spi"
	"github.com/embeddedgo/pico/hal/spi/spi0"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
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
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	// Setup SPI0 driver
	sm := spi0dma.Master()
	sm.UsePin(miso, spi.RXD)
	sm.UsePin(mosi, spi.TXD)
	sm.UsePin(csn, spi.CSN)
	sm.UsePin(sck, spi.SCK)
	sm.SetBaudrate(1e6)

	// Data to sent.
	s8 := ">> 0123456789 abcdefghijklmnoprstuvwxyz ABCDEFGHIJKLMNOPRSTUVWXYZ <<"
	s16 := make([]uint16, 77)
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
			fmt.Print("WriteStringRead ok\n")
		} else {
			fmt.Printf("WriteStringRead err: '%s'\n", buf8[:n])
		}

		for i := range 0x100 {
			b := byte(i)
			if sm.WriteReadByte(b) != b {
				fmt.Printf("WriteReadByte err: %x\n", b)
				goto ok1
			}
		}
		fmt.Print("WriteReadByte ok\n")
	ok1:

		sm.SetConfig(spi.Word16b)
		n = sm.WriteRead16(s16, buf16)
		if slices.Equal(s16, buf16[:n]) {
			fmt.Print("WriteRead16 ok\n")
		} else {
			fmt.Printf("WriteRead16 err: %x\n", buf16[:n])
		}

		for i := range 10000 {
			w := uint16(i)
			if sm.WriteReadWord16(w) != w {
				fmt.Printf("WriteReadWord16 err: %x\n", w)
				goto ok2
			}
		}
		fmt.Print("WriteReadWord16 ok\n")
	ok2:

		clear(buf8)
		clear(buf16)
		time.Sleep(time.Second)
	}
}
