// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import "github.com/embeddedgo/pico/hal/iomux"

type Signal int8

// Do not reorder the signal constants below. They must match the order in pins.

const (
	RXD Signal = iota
	CSN
	SCK
	TXD
)

// Pins returns the IO pins that can be used for the singal sig.
func (p *Periph) Pins(sig Signal) []iomux.Pin {
	n := len(rxPins) / 2
	rx := rxPins[n*num(p):][:n]
	pins := make([]iomux.Pin, n)
	for i := 0; i < len(rx); i++ {
		pins[i] = iomux.Pin(int(rx[i]) + int(sig))
	}
	return pins
}

// UsePin is a helper function that can be used to configure IO pins as required
// by the SPI peripheral. Only certain pins can be used (see datasheet). UsePin
// returns true on succes or false if it isn't possible to use a pin as a sig.
// See also Periph.Pins.
func (d *Master) UsePin(pin iomux.Pin, sig Signal) bool {
	n := len(rxPins) / 2
	rx := rxPins[n*num(d.p):][:n]
	rp := int(pin) - int(sig) // pin transformed to RX pin
	for i := 0; i < len(rx); i++ {
		if int(rx[i]) == rp {
			goto ok
		}
	}
	return false
ok:
	pin.SetAltFunc(pin.AltFunc()&^iomux.Func | iomux.SPI)
	switch sig {
	case RXD:
		pin.Setup(iomux.InpEn | iomux.OutDis)
	default: // CSN, SCK, TXD
		pin.Setup(iomux.D4mA)
	}
	return true
}

const rxPins = "" +
	p00 + p04 + p16 + p20 + p32 + p36 + // SPI0 RXD
	p08 + p12 + p24 + p28 + p40 + p44 //   SPI1 RXD

const (
	p00 = "\x00"
	p01 = "\x01"
	p02 = "\x02"
	p03 = "\x03"
	p04 = "\x04"
	p05 = "\x05"
	p06 = "\x06"
	p07 = "\x07"
	p08 = "\x08"
	p09 = "\x09"
	p10 = "\x0a"
	p11 = "\x0b"
	p12 = "\x0c"
	p13 = "\x0d"
	p14 = "\x0e"
	p15 = "\x0f"
	p16 = "\x10"
	p17 = "\x11"
	p18 = "\x12"
	p19 = "\x13"
	p20 = "\x14"
	p21 = "\x15"
	p22 = "\x16"
	p23 = "\x17"
	p24 = "\x18"
	p25 = "\x19"
	p26 = "\x1a"
	p27 = "\x1b"
	p28 = "\x1c"
	p29 = "\x1d"
	p30 = "\x1e"
	p31 = "\x1f"
	p32 = "\x20"
	p33 = "\x21"
	p34 = "\x22"
	p35 = "\x23"
	p36 = "\x24"
	p37 = "\x25"
	p38 = "\x26"
	p39 = "\x27"
	p40 = "\x28"
	p41 = "\x29"
	p42 = "\x2a"
	p43 = "\x2b"
	p44 = "\x2c"
	p45 = "\x2d"
	p46 = "\x2e"
	p47 = "\x2f"
)
