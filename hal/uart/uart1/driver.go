// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart1

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
		driver = uart.NewDriver(uart.UART(1))
		irq.UART1.Enable(rtos.IntPrioLow, system.NextCPU())
	}
	return driver
}

//go:interrupthandler
func _UART1_Handler() { driver.ISR() }

//go:linkname _UART1_Handler IRQ34_Handler
