// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dmairq allows to share the limited number of system-level DMA
// interrupts between the interrupt service routines (ISRs, interrupt hadlers)
// for individual DMA channels.
//
// The RP2350 has one DMA controller, 16 DMA channels and 4 system-level DMA
// intrrupts. In case of RP2350 this package uses only two first system-level
// DMA intrrupts, routing the interrupt 0 to CPU0 and the interrupt 1 to CPU1.
// For example, if you want to handle IRQs for some DMA channel on the CPU0 you
// should register this ISR for this channel using the SetISR function and next
// enable the interrupt 0 on this channel. If you want evenly distribute IRQs
// between two available CPUs use the value of int(system.NextCPU() & 1) as the
// interrupt number. In case of RP2350 the remaining two system-level DMA
// interrupts not used by this package can be used for channels and their ISRs
// in this specific cases for which the way this package handles/routes ISRs
// isn't suitable.
package dmairq

import (
	"embedded/rtos"

	"github.com/embeddedgo/pico/hal/dma"
)

// SetISR sets the isr function to be an interupt handler for the channel c. The
// isr function itself and all functions it calls sholud have go:nosplit
// directive (go:nowritebarrierrec is also recommended).
func SetISR(c dma.Channel, isr func()) { setISR(c, isr) }

func init() { enableIRQs(rtos.IntPrioLow) }
