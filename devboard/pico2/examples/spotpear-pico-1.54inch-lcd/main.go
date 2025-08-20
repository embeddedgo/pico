// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pico-lcd-1.54inch-lcd allow to test the onboar buttons.
package main

import (
	"image"
	"image/color"

	"github.com/embeddedgo/display/pix"
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/devboard/pico2/module/spotpear/pico-1.54inch-lcd/buttons"
	"github.com/embeddedgo/pico/devboard/pico2/module/spotpear/pico-1.54inch-lcd/lcd"
)

func main() {
	disp := lcd.Display
	a := disp.NewArea(disp.Bounds())
	a.SetColor(color.Black)
	a.Fill(a.Bounds())

	xh := (a.Bounds().Min.X + a.Bounds().Max.X) / 2
	buttons := []*Square{
		NewSquare(xh-44, 20, rgb(255, 0, 0), buttons.L),
		NewSquare(xh+44, 20, rgb(0, 255, 0), buttons.R),
		NewSquare(xh-44, 50, rgb(0, 0, 255), buttons.Select),
		NewSquare(xh+44, 50, rgb(255, 255, 255), buttons.Start),
		NewSquare(xh-60, 90, rgb(0, 255, 255), buttons.Up),
		NewSquare(xh+60, 90, rgb(255, 255, 0), buttons.Y),
		NewSquare(xh-40, 130, rgb(255, 0, 255), buttons.Right),
		NewSquare(xh+40, 130, rgb(255, 0, 0), buttons.X),
		NewSquare(xh-80, 130, rgb(0, 255, 0), buttons.Left),
		NewSquare(xh+80, 130, rgb(255, 255, 255), buttons.B),
		NewSquare(xh-60, 170, rgb(0, 255, 255), buttons.Down),
		NewSquare(xh+60, 170, rgb(255, 255, 0), buttons.A),
	}

	for {
		for _, b := range buttons {
			b.Draw(a)
		}
	}
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
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
	a := image.Point{s.X - 12, s.Y - 12}
	b := image.Point{s.X + 12, s.Y + 12}
	state := s.B.Read()
	if s.last != state {
		s.last = state
		ar.SetColor(color.Black)
		ar.RoundRect(a, b, 0, 0, true)
		ar.SetColor(s.C)
		ar.RoundRect(a, b, 0, 0, state != 0)
	}
}
