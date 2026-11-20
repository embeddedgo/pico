// Copyright 2026 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Program display shows how to use the onnboard OLED display.
package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/embeddedgo/display/font/subfont/font9/terminus12"
	"github.com/embeddedgo/pico/devboard/gnpe042/board/oled"
)

var (
	black = color.Gray{0}
	white = color.Gray{255}
)

// The imported oled package provides the initialized pix.Display.
//
// The visible part of the GDDRAM for this 70x40 monochrome display is 360
// bytes. Even in case of the full sized 128x64 displays based on the same
// controller the most efficient way to handle drawing is to use a framebuffer
// in SRAM and we, of course, use it here to. You need to use the Area.Flush or
// Display.Flush methods to see the latest changes in the frame buffer.

func main() {
	// Set the drawing area for the whole display and clear it.
	a := oled.Display.NewArea(oled.Display.Bounds())
	a.SetColor(black)
	a.Fill(a.Bounds())

	// Draw a rounded frame around the screen.
	r := a.Bounds().Inset(4)
	a.SetColor(white)
	a.RoundRect(r.Min, r.Max.Sub(image.Pt(1, 1)), 4, 4, false)

	// Blink the center of the display 4 times.
	r = a.Bounds().Inset(8)
	for i := 0; i < 6; i++ {
		a.SetColor(color.Gray{255 * uint8(i&1)})
		a.Fill(r)
		a.Flush()
		time.Sleep(time.Second)
	}

	// Setup a TextWriter on the drawing area.
	ff := terminus12.NewFace(terminus12.X0020_007e)
	w := a.NewTextWriter(ff)
	w.SetColor(white) // text color
	a.SetColor(black)

	for i := uint64(0); ; i++ {
		a.Fill(a.Bounds())     // clear the whole screen
		w.Pos = a.Bounds().Min // reset the text writer position
		fmt.Fprint(w, i)       // print a number
		a.Flush()
	}
}
