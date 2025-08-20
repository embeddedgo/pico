// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"embedded/rtos"
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/embeddedgo/device/bus/i2cbus"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/iomux"
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

	sda.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	sda.SetAltFunc(iomux.I2C)
	scl.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	scl.SetAltFunc(iomux.I2C)

	d = i2c.NewMaster(i2c.I2C(0))
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
		State.Store(1)
		d.Flush()
		State.Store(2)
		d.Wait(i2c.TX_EMPTY)
		State.Store(3)
		d.Abort() // stop

		// Wait for the end of write
		State.Store(4)
		for {
			d.WriteBytes(addr)
			d.Flush()
			State.Store(5)
			d.Wait(i2c.TX_EMPTY)
			State.Store(6)
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
		State.Store(7)
		d.ReadBytes(in[:n])
		State.Store(8)
		d.Abort() // stop
		State.Store(9)

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

var State atomic.Int32

//go:interrupthandler
func I2C0_Handler() {
	d.ISR()
}
