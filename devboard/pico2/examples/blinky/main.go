// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
)

func main() {
	for {
		leds.User.SetOn()
		time.Sleep(100 * time.Millisecond)
		leds.User.SetOff()
		time.Sleep(400 * time.Millisecond)
	}
}
