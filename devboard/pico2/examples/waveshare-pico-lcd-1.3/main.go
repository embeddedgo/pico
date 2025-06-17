// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pico-lcd-1.3 is similar to the Waveshare demo code for this module. It allows
// you to test the display and the onboard joystick and buttons.
package main

import (
	"image"
	"image/color"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/devboard/pico2/module/waveshare/pico-lcd-1.3/buttons"
	"github.com/embeddedgo/pico/devboard/pico2/module/waveshare/pico-lcd-1.3/lcd"
)

func main() {
	buttons := []interface{ Draw(*pix.Area) }{
		// Joystick
		NewSquare(60, 120, color.White, buttons.Ctrl),
		NewTriangle(60, 80, color.White, 0, buttons.Up),
		NewTriangle(60, 160, color.White, 1, buttons.Down),
		NewTriangle(100, 120, color.White, 2, buttons.Right),
		NewTriangle(20, 120, color.White, 3, buttons.Left),

		// Buttons
		NewSquare(200, 30, color.RGBA{255, 0, 0, 255}, buttons.A),
		NewSquare(200, 90, color.RGBA{0, 255, 0, 255}, buttons.B),
		NewSquare(200, 150, color.RGBA{0, 0, 255, 255}, buttons.X),
		NewSquare(200, 210, color.RGBA{255, 255, 0, 255}, buttons.Y),
	}

	disp := lcd.Display
	a := disp.NewArea(disp.Bounds())
	a.SetColor(color.Black)
	a.Fill(a.Bounds())
	for {
		for _, b := range buttons {
			b.Draw(a)
		}
	}
}

type Square struct {
	X, Y int
	C    color.Color
	B    common.Button

	last int
}

func NewSquare(x, y int, c color.Color, b common.Button) *Square {
	t := &Square{X: x, Y: y, C: c, B: b}
	t.last = -1
	return t
}

func (s *Square) Draw(ar *pix.Area) {
	a := image.Point{s.X - 15, s.Y - 15}
	b := image.Point{s.X + 15, s.Y + 15}
	state := s.B.Read()
	if s.last != state {
		s.last = state
		ar.SetColor(color.Black)
		ar.RoundRect(a, b, 0, 0, true)
		ar.SetColor(s.C)
		ar.RoundRect(a, b, 0, 0, state != 0)
	}
}

type Triangle struct {
	X, Y int
	C    color.Color
	Dir  int
	B    common.Button

	last int
}

func NewTriangle(x, y int, c color.Color, dir int, b common.Button) *Triangle {
	t := &Triangle{X: x, Y: y, C: c, Dir: dir, B: b}
	t.last = -1
	return t
}

func (t *Triangle) Draw(ar *pix.Area) {
	var a, b, c image.Point
	switch t.Dir {
	case 0:
		a.X, a.Y = t.X-15, t.Y+13
		b.X, b.Y = t.X, t.Y-14
		c.X, c.Y = t.X+15, t.Y+13
	case 1:
		a.X, a.Y = t.X-15, t.Y-13
		b.X, b.Y = t.X, t.Y+14
		c.X, c.Y = t.X+15, t.Y-13
	case 2:
		a.X, a.Y = t.X-15, t.Y-14
		b.X, b.Y = t.X+15, t.Y
		c.X, c.Y = t.X-15, t.Y+14
	case 3:
		a.X, a.Y = t.X+15, t.Y-14
		b.X, b.Y = t.X-15, t.Y
		c.X, c.Y = t.X+15, t.Y+14
	}
	state := t.B.Read()
	if t.last != state {
		t.last = state
		ar.SetColor(color.Black)
		ar.Quad(a, b, c, a, true)
		ar.SetColor(t.C)
		ar.Quad(a, b, c, a, state != 0)
	}
}
