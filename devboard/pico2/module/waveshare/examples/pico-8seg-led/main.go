// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/devboard/pico2/module/waveshare/pico-8seg-led/segled"
)

const (
	fps    = 60
	digits = 4
	repeat = 100
)

var symbols = [8]byte{1, 2, 4, 8, 16, 32, 64, 128}

func main() {
	sr := segled.NewShiftReg(pins.GP11, pins.GP10, pins.GP9)
	for i := 0; ; i++ {
		for range repeat {
			for d := 0; d < digits; d++ {
				sr.WriteByte(^byte(1 << uint(d)))
				sr.WriteByte(symbols[(i+d)&7])
				sr.Latch()
				time.Sleep(time.Second / (fps * digits))
			}
		}
	}
}
