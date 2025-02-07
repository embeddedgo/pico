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
	"github.com/embeddedgo/pico/hal/uart"
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

	u := uart.UART(0)
	d := uart.NewDriver(u)
	d.Setup(uart.Word8b, 115200)
	u.CR.Store(uart.UARTEN | uart.TXE | uart.RXE)

	baudStr := strconv.Itoa(d.Baudrate())
	for i := 0; ; i++ {
		puts(u, fmt.Sprintf("i=%d baudrate=%s\n\r", i, baudStr))
		leds.User.Toggle()
	}
}
