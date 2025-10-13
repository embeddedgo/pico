// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

// Eeprom writes and reads the memory of the 24C64/128/256 I2C EEPROM. The only
// difference to the less dense 24C0x EEPROMs is the use of 16 bit address
// instead of 8 bit one. Check also../eepromll to see a similar example that
// uses the RP2350 specific low-level I2C interface.
import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/embeddedgo/device/bus/i2cbus"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c0"
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

	d := i2c0.Master()
	d.UsePin(sda, i2c.SDA)
	d.UsePin(scl, i2c.SCL)
	d.Setup(100e3)

	c := d.NewConn(prefix<<3 | a2a1a0)

	var out, in [pageSize]byte

loop:
	for page := 0; ; page++ {
		a := page * pageSize
		addr := []byte{byte(a >> 8), byte(a)} // 16-bit memory address

		n := rand.Intn(pageSize) + 1
		randomData(out[:n])

		fmt.Printf("Wr %2d B page %d: %s ", n, page, out[:n])
		c.Write(addr)
		c.Write(out[:n])
		err := c.Close()
		if err != nil {
			fmt.Println("\nWr error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Wait for the end of write
		for {
			c.Write(addr)
			err = c.Close()
			if err == nil {
				break
			}
			if !errors.Is(err, i2cbus.ErrACK) {
				fmt.Print("\nwait error: ", err, "\n")
				time.Sleep(2 * time.Second)
				continue loop
			}
			fmt.Print(".")
		}
		fmt.Println(" done")

		c.Read(in[:n])
		err = c.Close()
		if err != nil {
			fmt.Println("Rd error:", err)
			time.Sleep(2 * time.Second)
		} else if string(in[:n]) != string(out[:n]) {
			fmt.Printf("Rd %2d B BAD! %d: %s\n\n", n, page, in[:n])
			time.Sleep(2 * time.Second)
		} else {
			fmt.Print("Rd OK\n")
		}
		for i := range in {
			in[i] = ':'
		}
	}
}
