// Copyright 2025 The Embedded Go Authors. All rights reserved.
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
	"github.com/embeddedgo/pico/hal/dma/dmairq"
)

type DMACopier struct {
	ch   dma.Channel
	done rtos.Note
}

func NewDMACopier(ch dma.Channel) *DMACopier {
	ch.SetConfig(dma.En|dma.S32b|dma.IncR|dma.IncW|dma.Always, ch)
	ch.ClearIRQ()
	ch.EnableIRQ(0)
	d := &DMACopier{ch: ch}
	dmairq.SetISR(ch, d.isr)
	return d
}

func (d *DMACopier) Copy(dst, src []uint32) int {
	n := min(len(dst), len(src))
	if n == 0 {
		return 0
	}
	d.done.Clear()
	ch := d.ch
	ch.SetReadAddr(unsafe.Pointer(&src[0]))
	ch.SetWriteAddr(unsafe.Pointer(&dst[0]))
	ch.SetTransCount(n, dma.Normal)
	ch.Trig()
	d.done.Sleep(-1)
	return n
}

//go:nosplit
func (d *DMACopier) isr() {
	d.ch.ClearIRQ()
	d.done.Wakeup()
}

func main() {
	n := 20000
	src := make([]uint32, n)
	dst := make([]uint32, n)

	for i := range src {
		src[i] = uint32(i)
	}

	dmacp := NewDMACopier(dma.DMA(0).AllocChannel())

	delay := time.Second // blink slow if transfer OK
	for i := range 10 {
		n := dmacp.Copy(dst, src[i:])
		if !slices.Equal(dst[:n], src[i:]) {
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
