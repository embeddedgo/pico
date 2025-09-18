// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c

import "github.com/embeddedgo/pico/hal/iomux"

type Signal int8

// Do not reorder the signal constants below.

const (
	SDA Signal = iota
	SCL
)

// I2C can use any pin. I2Cn SDA: 0 + 2*n + 4*i, I2Cn SCL: 1 + 2*n + 4*i
const numPins = 48

// Pins returns the IO pins that can be used for the singal sig.
func (p *Periph) Pins(sig Signal) []iomux.Pin {
	pins := make([]iomux.Pin, numPins/4)
	for i := range pins {
		pins[i] = iomux.Pin(int(sig) + num(p)*2 + i*4)
	}
	return pins
}

// UsePin is a helper function that can be used to configure IO pins as required
// by the SPI peripheral. Only certain pins can be used (see datasheet). UsePin
// returns true on succes or false if it isn't possible to use a pin as a sig.
// See also Periph.Pins.
func (d *Master) UsePin(pin iomux.Pin, sig Signal) bool {
	if int(pin)&3 != int(sig)+num(d.p)*2 {
		return false
	}
	pin.SetAltFunc(pin.AltFunc()&^iomux.Func | iomux.I2C)
	pin.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	return true
}
