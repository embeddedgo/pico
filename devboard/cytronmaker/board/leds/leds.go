// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leds

import (
	"github.com/embeddedgo/pico/devboard/common"
	_ "github.com/embeddedgo/pico/devboard/cytronmaker/board/system"
	"github.com/embeddedgo/pico/hal/iomux"
)

// The LED on the Pico 2 board.
const User = common.LED(iomux.P25)

func init() {
	common.ConnectLED(User.Pin(), iomux.D8mA, 0)
}
