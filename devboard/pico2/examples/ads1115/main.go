// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/embeddedgo/device/adc/ads111x"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

const (
	addr = 0b100_1000 // address if the ADDR pin is connected to GND
	cfg  = ads111x.OS | ads111x.AIN0_AIN1 | ads111x.FS2048 | ads111x.SINGLESHOT | ads111x.R8
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
	sda.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	sda.SetAltFunc(iomux.I2C)
	scl.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	scl.SetAltFunc(iomux.I2C)

	d := i2c.NewMaster(i2c.I2C(0))
	d.Setup(100e3)
	d.SetAddr(addr)

again:
	// Start a new A/D conversion.
	d.WriteCmds([]int16{
		i2c.Send | ads111x.RegCfg,
		i2c.Send | int16(cfg>>8),
		i2c.Send | int16(cfg&0xff),
	})

	var buf [2]byte

	// Wait for the end of conversion.
	for {
		d.WriteCmd(i2c.Recv | i2c.Restart)
		d.WriteCmd(i2c.Recv)
		d.ReadBytes(buf[:])
		if buf[0]&byte(ads111x.OS>>8) != 0 {
			break
		}
	}

	// Read the result of ADC.
	d.WriteCmds([]int16{
		i2c.Send | ads111x.RegV,
		i2c.Recv,
		i2c.Recv | i2c.Stop,
	})
	d.ReadBytes(buf[:])

	// Convert to volts and print.
	fmt.Printf("%.6f V\n", volt(buf))

	goto again
}

func volt(buf [2]byte) float64 {
	scale := 6.144
	if shift := cfg & ads111x.PGA >> ads111x.PGAn; shift != 0 {
		scale = 4.096 / float64(uint(1)<<(shift-1))
	}
	return float64(int16(buf[0])<<8|int16(buf[1])) * scale / 0x8000
}
