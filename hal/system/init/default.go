// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !nosysinit && !133MHz && !150MHz

package init

import (
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/system/timer/riscvst"
)

func init() {
	system.SetupPico2_125MHz()
	riscvst.Setup()
}
