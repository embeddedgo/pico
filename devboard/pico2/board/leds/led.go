// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leds

import (
	_ "github.com/embeddedgo/pico/devboard/pico2/board/system"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

const User LED = 25 // The onboard LED

type LED uint8

func (d LED) SetOn()     { gpio.BitForPin(iomux.Pin(d)).Set() }
func (d LED) SetOff()    { gpio.BitForPin(iomux.Pin(d)).Clear() }
func (d LED) Toggle()    { gpio.BitForPin(iomux.Pin(d)).Toggle() }
func (d LED) Set(on int) { gpio.BitForPin(iomux.Pin(d)).Store(on) }
func (d LED) Get() int   { return gpio.BitForPin(iomux.Pin(d)).LoadOut() }

func init() {
	pin := iomux.Pin(User)
	pin.Setup(iomux.D4mA)
	gpio.UsePin(pin)
	gpio.BitForPin(pin).EnableOut()
}
