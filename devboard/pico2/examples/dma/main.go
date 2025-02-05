// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dma shows how to use the DMA controller for RAM to RAM transfers.
package main

import (
	"embedded/rtos"
	"slices"
	"time"
	"unsafe"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/irq"
)

var (
	ch  dma.Channel
	tce rtos.Note
)

func main() {
	n := 20000
	src := make([]uint32, n)
	dst := make([]uint32, n)
	srcAddr := unsafe.Pointer(&src[0])
	dstAddr := unsafe.Pointer(&dst[0])

	for i := range src {
		src[i] = uint32(i)
	}

	irq.DMA_IRQ_0.Enable(rtos.IntPrioLow, 1) // enable DMA IRQ 0 on Core1

	ch = dma.DMA(0).AllocChannel()
	ch.SetTransCount(n, dma.Normal)
	ch.ClearIRQ()
	ch.EnableIRQ(0)
	ch.SetConf(dma.En|dma.S32b|dma.IncR|dma.IncW|dma.Always, ch)

	delay := time.Second // blink slow if transfer OK

	for range 10 {
		tce.Clear()
		ch.SetReadAddr(srcAddr)
		ch.SetWriteAddr(dstAddr)
		ch.Trig()
		tce.Sleep(-1)

		if !slices.Equal(src, dst) {
			delay /= 8 // blink fast in case of transfer error
			break
		}
		clear(dst)
	}

	for {
		leds.User.SetOn()
		time.Sleep(delay)
		leds.User.SetOff()
		time.Sleep(delay)
	}
}

//go:interrupthandler
func DMA_IRQ_0_Handler() {
	ch.ClearIRQ()
	tce.Wakeup()
}
