// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"embedded/rtos"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/embeddedgo/device/bus/i2cbus"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

const (
	prefix   = 0b1010 // address prefix (0xa)
	a2a1a0   = 0b000  // address pins
	pageSize = 32
)

func randomData(p []byte) {
	for i := range p {
		p[i] = byte(rand.Int31()&63 + ';')
	}
}

var d *i2c.Master

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		sda   = pins.GP20
		scl   = pins.GP21
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	d = i2c.NewMaster(i2c.I2C(0))
	d.UsePin(sda, i2c.SDA)
	d.UsePin(scl, i2c.SCL)
	d.Setup(100e3)
	irq.I2C0.Enable(rtos.IntPrioLow, 0)
	d.SetAddr(prefix<<3 | a2a1a0)

	var out, in [pageSize]byte

loop:
	for page := 0; ; page++ {
		a := page * pageSize
		addr := []byte{byte(a >> 8), byte(a)}

		n := rand.Intn(pageSize) + 1
		randomData(out[:n])

		fmt.Printf("Wr %2d B page %d: %s ", n, page, out[:n])
		d.WriteBytes(addr)
		d.WriteBytes(out[:n])
		d.Flush()
		d.Wait(i2c.TX_EMPTY)
		d.Abort() // stop

		// Wait for the end of write
		for {
			d.WriteBytes(addr)
			d.Flush()
			d.Wait(i2c.TX_EMPTY)
			err := d.Err(true)
			if err == nil {
				break
			}
			if !errors.Is(err, i2cbus.ErrACK) {
				fmt.Print("\n", err, "\n")
				time.Sleep(time.Second)
				continue loop
			}
			fmt.Print(".")
		}
		fmt.Println(" done")

		d.WriteCmd(i2c.Recv | int16(n-1))
		d.ReadBytes(in[:n])
		d.Abort() // stop

		err := d.Err(true)
		if err != nil {
			fmt.Println(err)
		} else if string(in[:n]) != string(out[:n]) {
			fmt.Printf("Rd %2d B BAD! %d: %s\n\n", n, page, in[:n])
		} else {
			fmt.Print("Rd OK\n")
		}
	}
}

//go:interrupthandler
func I2C0_Handler() {
	d.ISR()
}
