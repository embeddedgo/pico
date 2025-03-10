// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leds

import (
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/iomux"
)

// The onboard LEDs
const (
	Blue  = common.LED(iomux.P24)
	Green = common.LED(iomux.P25)
	User  = Green
)

func init() {
	common.ConnectLED(Blue.Pin(), iomux.D8mA, 0)
	common.ConnectLED(Green.Pin(), iomux.D8mA, 0)
}
