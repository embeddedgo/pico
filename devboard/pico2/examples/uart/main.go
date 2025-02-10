// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"fmt"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
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

var d *uart.Driver

func main() {
	tx := pins.GP0
	rx := pins.GP1

	tx.Setup(iomux.D2mA)
	rx.Setup(iomux.InpEn)
	tx.SetAltFunc(iomux.UART)
	rx.SetAltFunc(iomux.UART)

	u := uart.UART(0)
	d = uart.NewDriver(u)
	d.Setup(uart.Word8b, 115200)
	d.EnableTx()
	irq.UART0.Enable(rtos.IntPrioLow, 0)
	s := "+ 0123456780 - abcdefghijklmnoprstuvwxyz - ABCDEFGHIJKLMNOPRSTUVWXYZ +\r\n"
	for {
		t0 := time.Now()
		const n = 256
		for i := 0; i < n; i++ {
			d.WriteString(s)
		}
		d.Flush()
		dt := time.Now().Sub(t0)
		baud := time.Duration(n*10*len(s)) * time.Second / dt
		fmt.Fprintf(d, "%v, %d baud\r\n", dt, baud)
		time.Sleep(5 * time.Second)
	}
}

//go:interrupthandler
func UART0_Handler() {
	d.ISR()
	leds.User.Toggle()
}
