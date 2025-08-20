// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Uart allows to tests the communication over UART peripheral.
package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

func main() {
	tx := pins.GP0
	rx := pins.GP1

	u := uart0.Driver()
	u.UsePin(tx, uart.TXD)
	u.UsePin(rx, uart.RXD)
	u.Setup(uart.Word8b, 115200)
	u.EnableTx()
	u.EnableRx()

	s := "+@#$%- 0123456780 - abcdefghijklmnoprstuvwxyz - ABCDEFGHIJKLMNOPRSTUVWXYZ -%$#@+\r\n"
	t0 := time.Now()
	const n = 100
	for i := 0; i < n; i++ {
		u.WriteString(s)
	}
	u.WaitTxDone()
	dt := time.Now().Sub(t0)
	speed := (time.Duration(n*10*len(s))*time.Second + dt/2) / dt
	baud := u.Baudrate()
	fmt.Fprintf(u, "\r\n%v, speed: %d baud, uart hw: %d baud\r\n", dt, speed, baud)

	var buf [128]byte

	for {
		u.WriteString("> ")
		n, err := u.Read(buf[:])
		u.Write(buf[:n])
		u.WriteString("\r\n")
		if err != nil {
			fmt.Fprintf(u, "error: %v\r\n", err)
		}
	}

}
