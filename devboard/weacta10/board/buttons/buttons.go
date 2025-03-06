// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buttons

import (
	"github.com/embeddedgo/pico/devboard/common"
	"github.com/embeddedgo/pico/hal/iomux"
	_ "github.com/embeddedgo/pico/hal/system/init"
)

// The onboard buttons.
const User = common.Button(iomux.P23)

func init() {
	common.ConnectButton(User.Pin(), 0, iomux.InpInvert)
}
