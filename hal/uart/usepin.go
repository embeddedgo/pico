// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import "github.com/embeddedgo/pico/hal/iomux"

type Signal int8

// Do not reorder the signal constants below. They must match the order in pins.

const (
	TXD Signal = iota
	RXD
	CTS
	RTS
)

// Pins returns the IO pins that can be used for the singal sig.
func (p *Periph) Pins(sig Signal) []iomux.Pin {
	n := len(txPins) / 2
	tx := txPins[n*num(p):][:n]
	if sig <= RXD {
		n *= 2
	}
	pins := make([]iomux.Pin, n)
	k := 0
	for i := 0; i < len(tx); i++ {
		pins[k] = iomux.Pin(int(tx[i]) + int(sig))
		k++
		if sig <= RXD {
			pins[k] = iomux.Pin(int(tx[i]) + int(sig) + 2)
			k++
		}
	}
	return pins
}

// UsePin is a helper function that can be used to configure IO pins as required
// by the UART peripheral. Only certain pins can be used (see datasheet). UsePin
// returns true on succes or false if it isn't possible to use a pin as a sig.
// See also Periph.Pins.
func (d *Driver) UsePin(pin iomux.Pin, sig Signal) bool {
	n := len(txPins) / 2
	tx := txPins[n*num(d.p):][:n]
	tp := int(pin) - int(sig) // pin transformed to TX pin

	af := iomux.UART
	for i := 0; i < len(tx); i++ {
		if int(tx[i]) == tp {
			goto ok
		}
		if sig <= RXD && int(tx[i]) == tp-2 {
			af = iomux.UART_AUX
			goto ok
		}
	}
	return false
ok:
	pin.SetAltFunc(pin.AltFunc()&^iomux.Func | af)
	switch sig {
	case TXD, RTS:
		pin.Setup(iomux.D4mA)
	default: // RXD, CTS
		pin.Setup(iomux.InpEn | iomux.OutDis)
	}
	return true
}

const txPins = "" +
	p00 + p12 + p16 + p28 + p32 + p44 + // UART0 TXD
	p04 + p08 + p20 + p24 + p36 + p40 //   UART1 TXD

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
