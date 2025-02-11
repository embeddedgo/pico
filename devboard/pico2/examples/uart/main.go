// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"fmt"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/uart"
)

var u *uart.Driver

func main() {
	tx := pins.GP0
	rx := pins.GP1

	tx.Setup(iomux.D2mA)
	rx.Setup(iomux.InpEn)
	tx.SetAltFunc(iomux.UART)
	rx.SetAltFunc(iomux.UART)

	u = uart.NewDriver(uart.UART(0))
	u.Setup(uart.Word8b, 230400)
	u.EnableTx()
	irq.UART0.Enable(rtos.IntPrioLow, 0)
	s := "+@#$%- 0123456780 - abcdefghijklmnoprstuvwxyz - ABCDEFGHIJKLMNOPRSTUVWXYZ -%$#@+\r\n"
	for {
		t0 := time.Now()
		const n = 256
		for i := 0; i < n; i++ {
			u.WriteString(s)
		}
		u.WaitTxDone()
		dt := time.Now().Sub(t0)
		baud := (time.Duration(n*10*len(s))*time.Second + dt/2) / dt
		fmt.Fprintf(u, "%v, %d baud\r\n", dt, baud)
		time.Sleep(3 * time.Second)
	}
}

//go:interrupthandler
func UART0_Handler() {
	u.ISR()
}
