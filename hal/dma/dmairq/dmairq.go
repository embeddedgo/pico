// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dmairq

import (
	"embedded/rtos"

	"github.com/embeddedgo/pico/hal/dma"
)

func SetISR(c dma.Channel, isr func()) { setISR(c, isr) }

func init() { enableIRQs(rtos.IntPrioLow) }
