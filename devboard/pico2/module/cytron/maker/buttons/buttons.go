// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buttons

import (
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/iomux"
)

// The onboard buttons.
const (
	B0 = common.Button(iomux.P20)
	B1 = common.Button(iomux.P21)
	B2 = common.Button(iomux.P22)
)

func init() {
	common.ConnectButton(B0.Pin(), 0, iomux.InpInvert)
	common.ConnectButton(B1.Pin(), 0, iomux.InpInvert)
	common.ConnectButton(B2.Pin(), 0, iomux.InpInvert)
}
