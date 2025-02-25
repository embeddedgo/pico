// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import "github.com/embeddedgo/pico/p/resets"

func BusyWaitAtLeastCycles(n uint)

func SetReset(pmask uint32, assert bool) {
	RESETS := resets.RESETS()
	if assert {
		AtomicSet(&RESETS.RESET, pmask)
	} else {
		AtomicClear(&RESETS.RESET, pmask)
		for RESETS.RESET_DONE.LoadBits(pmask) != pmask {
		}
	}
}
