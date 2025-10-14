// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c0dma

import (
	"embedded/rtos"
	_ "unsafe"

	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/internal"
	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/system"
)

var master *i2c.Master

// Master returns ready to use driver for I2C master.
func Master() *i2c.Master {
	if master == nil {
		master = internal.MasterDMA(i2c.I2C(0))
		irq.I2C0.Enable(rtos.IntPrioLow, system.NextCPU())
	}
	return master
}

//go:interrupthandler
func _I2C0_Handler() { master.ISR() }

//go:linkname _I2C0_Handler IRQ36_Handler
