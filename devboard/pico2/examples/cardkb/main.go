// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c0dma"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		sda   = pins.GP8
		scl   = pins.GP9
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	fmt.Print("\nWelcome!\n\n")

	// I2C0
	i2cm := i2c0dma.Master()
	i2cm.UsePin(sda, i2c.SDA)
	i2cm.UsePin(scl, i2c.SCL)
	i2cm.Setup(100e3)

	// M5Stack CardKB I2C keyboard
	kb := i2cm.NewConn(0x5F)

	var (
		buf [1]byte
		err error
	)
	for {
		buf[0], err = kb.ReadByte()
		kb.Close()
		if err != nil {
			fmt.Println("error:", err)
			time.Sleep(time.Second)
			continue
		}
		if buf[0] != 0 {
			os.Stdout.Write(buf[:])
		}
	}

}
