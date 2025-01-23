// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system/clock"
	"github.com/embeddedgo/pico/p/resets"
	"github.com/embeddedgo/pico/p/uart"
)

func putc(u *uart.Periph, c byte) {
	for u.FR.LoadBits(uart.TXFF) != 0 {
	}
	u.DR.Store(uint32(c))
}

func puts(u *uart.Periph, s string) {
	for i := 0; i < len(s); i++ {
		putc(u, s[i])
	}
}

func main() {
	tx := pins.GP0
	rx := pins.GP1

	tx.Setup(iomux.D2mA)
	rx.Setup(iomux.InpEn)
	tx.SetAltFunc(iomux.UART)
	rx.SetAltFunc(iomux.UART)

	u := uart.UART0()

	RESETS := resets.RESETS()
	RESETS.RESET.ClearBits(resets.UART0)
	for RESETS.RESET_DONE.LoadBits(resets.UART0) == 0 {
	}

	baud := 115200
	periHz := clock.PERI.Freq()
	brdiv := uint32(8*periHz/int64(baud)) + 1
	ibrd := brdiv >> 7
	fbrd := brdiv & 0x7f >> 1
	if ibrd == 0 {
		ibrd = 1
		fbrd = 0
	} else if ibrd >= 0xffff {
		ibrd = 0xffff
		fbrd = 0
	}
	u.IBRD.Store(ibrd)
	u.FBRD.Store(fbrd)
	u.LCR_H.Store(u.LCR_H.Load()) // dummy write to latch IBRD and FBRD

	u.LCR_H.Store(3<<uart.WLENn | uart.FEN)
	u.CR.Store(uart.UARTEN | uart.TXE | uart.RXE)

	baudStr := strconv.Itoa(int(4 * periHz / int64(ibrd<<6+fbrd)))
	for i := 0; ; i++ {
		puts(u, fmt.Sprintf("%d: baudrate=%s\n\r", i, baudStr))
		leds.User.Toggle()
	}
}
