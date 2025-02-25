// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart0

import (
	"embedded/rtos"
	_ "unsafe"

	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/uart"
)

var driver *uart.Driver

// Driver returns ready to use driver for UART.
func Driver() *uart.Driver {
	if driver == nil {
		driver = uart.NewDriver(uart.UART(0))
		irq.UART0.Enable(rtos.IntPrioLow, system.NextCPU())
	}
	return driver
}

//go:interrupthandler
func _UART0_Handler() { driver.ISR() }

//go:linkname _UART0_Handler IRQ33_Handler
