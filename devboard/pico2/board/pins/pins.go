// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pins

import (
	_ "github.com/embeddedgo/pico/devboard/pico2/board/system"
	"github.com/embeddedgo/pico/hal/iomux"
)

const (
	GP0     = iomux.P00
	GP1     = iomux.P01
	GP2     = iomux.P02
	GP3     = iomux.P03
	GP4     = iomux.P04
	GP5     = iomux.P05
	GP6     = iomux.P06
	GP7     = iomux.P07
	GP8     = iomux.P08
	GP9     = iomux.P09
	GP10    = iomux.P10
	GP11    = iomux.P11
	GP12    = iomux.P12
	GP13    = iomux.P13
	GP14    = iomux.P14
	GP15    = iomux.P15
	GP16    = iomux.P16
	GP17    = iomux.P17
	GP18    = iomux.P18
	GP19    = iomux.P19
	GP20    = iomux.P20
	GP21    = iomux.P21
	GP22    = iomux.P22
	GP26_A0 = iomux.P26
	GP27_A1 = iomux.P27
	GP28_A2 = iomux.P28
)

const (
	PS   = iomux.P23 // RT6150 Power Save pin
	VBUS = iomux.P24 // VBUS presence
	VSYS = iomux.P29 // analog VSYS/3
)
