// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buttons

import (
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/iomux"
	_ "github.com/embeddedgo/pico/hal/system/init"
)

const (
	// Joystick
	Up    = common.Button(iomux.P02)
	Down  = common.Button(iomux.P18)
	Left  = common.Button(iomux.P16)
	Right = common.Button(iomux.P20)
	Ctrl  = common.Button(iomux.P03)

	// Buttons
	A = common.Button(iomux.P15)
	B = common.Button(iomux.P17)
	X = common.Button(iomux.P19)
	Y = common.Button(iomux.P21)
)

func init() {
	for _, btn := range []common.Button{Up, Down, Left, Right, Ctrl, A, B, X, Y} {
		common.ConnectButton(btn.Pin(), iomux.PullUp, iomux.InpInvert)
	}
}
