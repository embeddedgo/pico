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
	"embedded/rtos"
	"os"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
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

var m *i2c.Master

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
	sda.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	sda.SetAltFunc(iomux.I2C)
	scl.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	scl.SetAltFunc(iomux.I2C)

	m = i2c.NewMaster(i2c.I2C(0))
	m.Setup(100e3)
	irq.I2C0.Enable(rtos.IntPrioLow, 0)

	m.SetAddr(0b010_0111)

	n := 1000
	cmds := make([]int16, n)
	for i := range cmds {
		cmds[i] = Z
	}
	cmds[n*0/4] = R
	cmds[n*1/4] = G
	cmds[n*2/4] = B
	cmds[n*3/4] = W
	cmds[n-1] = Z

	for {
		m.WriteCmds(cmds)
		if err := m.Err(true); err != nil {
			os.Stderr.WriteString("\n")
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
			time.Sleep(time.Second)
		} else {
			//os.Stderr.WriteString(".")
		}
	}
}

//go:interrupthandler
func I2C0_Handler() {
	m.ISR()
}
