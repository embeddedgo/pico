// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embedded/rtos"
	"time"

	"github.com/embeddedgo/led"
	"github.com/embeddedgo/led/ws281x/wsuart"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/uart"
)

var u *uart.Driver

func main() {
	tx := pins.GP22

	tx.Setup(iomux.D4mA)
	tx.SetAltFunc(iomux.UART_AUX | iomux.OutInvert)

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	u = uart.NewDriver(uart.UART(1))
	u.Setup(uart.Word7b, 3000000000/1390)
	u.EnableTx()
	irq.UART1.Enable(rtos.IntPrioLow, 0)

	grb := wsuart.GRB
	strip := wsuart.Make(8 * 8)

	for {
		for _, c := range []led.Color{
			led.RGB(255, 0, 0),
			led.RGB(0, 255, 0),
			led.RGB(0, 0, 255),
			led.RGB(255, 255, 0),
			led.RGB(0, 255, 255),
			led.RGB(255, 0, 255),
			led.RGB(255, 255, 255),
		} {
			pixel := grb.Pixel(c)
			for i := 0; i < 64; i += 8 {
				strip.Clear()
				for k := i; k < i+8; k++ {
					strip[k] = pixel
				}
				u.Write(strip.Bytes())
				time.Sleep(time.Second / 2)
			}
		}
	}
}

//go:interrupthandler
func UART1_Handler() {
	u.ISR()
}
