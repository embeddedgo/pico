// Copyright 2022 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build rp2350

package dmairq

import (
	"math/bits"
	"sync/atomic"
	"unsafe"

	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/irq"
)

var handlers [16]unsafe.Pointer // func()

func setISR(c dma.Channel, isr func()) {
	h := *(*unsafe.Pointer)(unsafe.Pointer(&isr))
	atomic.StorePointer(&handlers[c.Num()], h)
}

//go:nosplit
func isr(irqn int) {
	chs := dma.DMA(0).ActiveIRQs(irqn)
	for i := 0; chs != 0; {
		n := bits.TrailingZeros32(chs)
		i += n
		chs >>= uint(n + 1)
		if h := atomic.LoadPointer(&handlers[i]); h != nil {
			(*(*func())(unsafe.Pointer(&h)))()
		}
	}
}

func enableIRQs(prio int) {
	irq.DMA0_0.Enable(prio, 0)
	irq.DMA0_1.Enable(prio, 1)
}

//go:interrupthandler
func _DMA0_0_Handler() { isr(0) }

//go:interrupthandler
func _DMA0_1_Handler() { isr(1) }

//go:linkname _DMA0_0_Handler IRQ10_Handler
//go:linkname _DMA0_1_Handler IRQ11_Handler
