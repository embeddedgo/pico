// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pico-restouch-lcd-2.8 allows to test the LCD module.
package main

import (
	"github.com/embeddedgo/display/pix/examples"
	"github.com/embeddedgo/pico/devboard/pico2/module/waveshare/pico-restouch-lcd-2.8/lcd"
)

func main() {
	disp := lcd.Display
	for {
		examples.RotateDisplay(disp)
		examples.DrawText(disp)
		examples.GraphicsTest(disp)
	}
}
