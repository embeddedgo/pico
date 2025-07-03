// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buttons

import (
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/iomux"
)

const (
	Left  = common.Button(iomux.P13)
	Down  = common.Button(iomux.P14)
	Up    = common.Button(iomux.P16)
	Right = common.Button(iomux.P17)

	X = common.Button(iomux.P02)
	A = common.Button(iomux.P03)
	Y = common.Button(iomux.P04)
	B = common.Button(iomux.P06)

	L = common.Button(iomux.P18)
	R = common.Button(iomux.P01)

	Select = common.Button(iomux.P26)
	Start  = common.Button(iomux.P27)
)

func init() {
	for _, btn := range []common.Button{Left, Down, Up, Right, X, A, Y, B, L, R, Select, Start} {
		common.ConnectButton(btn.Pin(), iomux.PullUp, iomux.InpInvert)
	}
}
