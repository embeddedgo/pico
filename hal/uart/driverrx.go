// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import "github.com/embeddedgo/pico/hal/internal"

func (d *Driver) EnableRx() {
	internal.AtomicSet(&d.p.CR, UARTEN|RXE)
}

func (d *Driver) DisableRx() {
	internal.AtomicClear(&d.p.CR, RXE)
}
