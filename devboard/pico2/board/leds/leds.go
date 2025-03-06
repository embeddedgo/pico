// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leds

import (
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/iomux"
	_ "github.com/embeddedgo/pico/hal/system/init"
)

// The onboard LED
const User = common.LED(iomux.P25)

func init() {
	common.ConnectLED(User.Pin(), iomux.D8mA, 0)
}
