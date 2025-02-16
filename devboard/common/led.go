// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

type LED uint8

func ConnectLED(pin iomux.Pin, drive iomux.Config, mod iomux.AltFunc) LED {
	pin.Setup(drive)
	pin.SetAltFunc(mod | iomux.GPIO)
	gpio.BitForPin(pin).EnableOut()
	return LED(pin)
}

//go:nosplit
func (d LED) SetOn() { gpio.BitForPin(iomux.Pin(d)).Set() }

//go:nosplit
func (d LED) SetOff() { gpio.BitForPin(iomux.Pin(d)).Clear() }

//go:nosplit
func (d LED) Toggle() { gpio.BitForPin(iomux.Pin(d)).Toggle() }

//go:nosplit
func (d LED) Set(on int) { gpio.BitForPin(iomux.Pin(d)).Store(on) }

//go:nosplit
func (d LED) Get() int { return gpio.BitForPin(iomux.Pin(d)).LoadOut() }

//go:nosplit
func (d LED) Pin() iomux.Pin { return iomux.Pin(d) }
