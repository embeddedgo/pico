// Copyright 2026 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"time"

	"github.com/embeddedgo/pico/devboard/gnpe042oled/board/oled"
)

func main() {
	a := oled.Display.NewArea(oled.Display.Bounds())
	r := a.Bounds().Inset(3)
	a.SetColorRGBA(255, 255, 255, 255)
	a.RoundRect(r.Min, r.Max.Sub(image.Pt(1, 1)), 3, 3, false)
	r = a.Bounds().Inset(8)

	for b := byte(0); ; b = ^b {
		a.SetColorRGBA(b, b, b, 255)
		a.Fill(r)
		a.Flush()
		time.Sleep(time.Second)

		//examples.RotateDisplay(oled.Display)
		//examples.DrawText(oled.Display)
		//examples.GraphicsTest(oled.Display)
	}
}
