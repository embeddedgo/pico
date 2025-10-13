// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pcf8574 writes consecutive numbers to the remote I/O expander chip (PCF8574)
// using I2C protocol.
//
// The easiest way to try this example is to use a PCF8574 based module intended
// for LCD displays and one or more LEDs. Low-voltage LEDs like red ones
// require a current limiting resistor of the order 150-200 Ω. High voltage LEDs
// like the white ones may work without any resistor.
//
// Connect your LEDs between pin 2 (closest to the I2C connector, 3.3V) and pins
// 4, 5, 6 (PCF8574 P0, P1, P2 outputs). If you have more LEDs you can connect
// four more to the pins 11, 12, 13, 14. Polarity matters. Pin 2 should be
// connected to the anodes of all LEDs. The easiest way to do it is to use a
// breadboard. Next connect the module pins GND, VCC, SDA, SCL to the Pico
// pins GND, 3.3V, 20, 21. After programming your Teensy with this example the
// LEDs should start blinking with different frequencies.
//
// As the LEDs are connected between 3.3V and P0, P1, P2 writing the
// corresponding bit zero turns the LED on, writting it one turns it off.
// Because of its quasi-bidirectional I/O the PCF8574 can't source enough
// current to stable light a LED connected between Px pin and GND (ones are
// weak so if you set pin high it can work as an input). On the other hand
// PCF8574 can sink enough current to light the LEDs (zeros are strong).
package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c0"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

const (
	R = 0b1111_1110
	G = 0b1111_1101
	B = 0b1111_1011
	W = 0b1110_1111
	Z = 0b1111_1111
)

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

	// I2C
	m := i2c0.Master()
	m.UsePin(sda, i2c.SDA)
	m.UsePin(scl, i2c.SCL)
	m.Setup(100e3)

	m.SetAddr(0b010_0111)

	// Demonstrate long I2C transfers
	n := 4000
	data := make([]byte, n)
	for i := range data {
		data[i] = Z
	}
	for i := 0; i < 20; i++ {
		data[n*0/4+i] = R
		data[n*1/4+i] = G
		data[n*2/4+i] = B
		data[n*3/4+i] = W
	}

	for {
		m.WriteBytes(data)
		if err := m.Err(true); err != nil {
			fmt.Println(err)
			for range 6 {
				leds.User.Toggle()
				time.Sleep(time.Second / 4)
			}
			time.Sleep(time.Second)
		}
	}
}
