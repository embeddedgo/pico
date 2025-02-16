// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

type Button uint8

func ConnectButton(pin iomux.Pin, pull iomux.Config, mod iomux.AltFunc) Button {
	pin.Setup(iomux.InpEn | iomux.OutDis | iomux.Schmitt | pull)
	pin.SetAltFunc(iomux.GPIO | mod)
	return Button(pin)
}

//go:nosplit
func (b Button) Read() int { return gpio.BitForPin(iomux.Pin(b)).Load() }

//go:nosplit
func (b Button) Pin() iomux.Pin { return iomux.Pin(b) }
