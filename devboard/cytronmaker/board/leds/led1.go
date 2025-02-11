// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// This implementation slightly worse inlines with go1.22 than the one from
// led.go.

package leds

import (
	_ "github.com/embeddedgo/pico/devboard/cytronmaker/board/system"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

// The onboard LED
var User = Connect(iomux.P25, iomux.D4mA, false)

type LED struct{ bit gpio.Bit }

func Connect(pin iomux.Pin, drive iomux.Config, invert bool) LED {
	pin.Setup(drive)
	af := iomux.GPIO
	if invert {
		af |= iomux.OutInvert
	}
	pin.SetAltFunc(af)
	bit := gpio.BitForPin(pin)
	bit.EnableOut()
	return LED{bit}
}

func (d LED) SetOn()     { d.bit.Set() }
func (d LED) SetOff()    { d.bit.Clear() }
func (d LED) Toggle()    { d.bit.Toggle() }
func (d LED) Set(on int) { d.bit.Store(on) }
func (d LED) Get() int   { return d.bit.LoadOut() }
