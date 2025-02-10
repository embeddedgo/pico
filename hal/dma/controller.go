// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dma provides interface to the DMA controller. The interface is based
// on two main types: Controller and Channel.
//
// Controller represents an instance of DMA controller (RP2350 provides only one
// controller numbered as 0). Each controller provides 16 channels.
//
// Channel represents a DMA channel. As all available channels have identical
// functionality this package doesn't allow you to select a specific channel.
// Instead you can use Controller.AllocChannel to allocate an unused one.
package dma

import (
	"embedded/mmio"
	"embedded/rtos"
	"runtime"
	"sync"
	"unsafe"

	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

type channel struct {
	readAddr   mmio.U32
	writeAddr  mmio.U32
	transCount mmio.U32
	ctrlTrig   mmio.U32

	ctrl1           mmio.U32
	readAddr1       mmio.U32
	writeAddr1      mmio.U32
	transCountTrig1 mmio.U32

	ctrl2          mmio.U32
	transCount2    mmio.U32
	readAddr2      mmio.U32
	writeAddrTrig2 mmio.U32

	ctrl3         mmio.U32
	readAddr3     mmio.U32
	transCount3   mmio.U32
	readAddrTrig3 mmio.U32
}

type irq struct {
	r, e, f, s mmio.U32
}

type mpu struct {
	bar mmio.U32
	lar mmio.U32
}

type chDbg struct {
	ctdReq mmio.U32
	tcr    mmio.U32
}

// A Controller represents a DMA controller.
type Controller struct {
	ch            [16]channel
	irq           [4]irq
	timer         [4]mmio.U32
	multiChanTrig mmio.U32
	sniffCtrl     mmio.U32
	sniffData     mmio.U32
	_             uint32
	fifoLevels    mmio.U32
	_             uint32
	nChannels     mmio.U32
	_             [5]uint32
	seccfgCh      [16]mmio.U32
	seccfgIRQ     [4]mmio.U32
	seccfgMisc    mmio.U32
	_             [11]uint32
	mpuCtrl       mmio.U32
	mpu           [8]mpu
	_             [175]uint32
	chDbg         [16]chDbg
}

func init() {
}

// DMA returns n-th controller (RP2350 suports onlu conrtoller 0).
func DMA(n int) *Controller {
	if n != 0 {
		panic("wrong DMA number")
	}
	return (*Controller)(unsafe.Pointer(mmap.DMA_BASE))
}

var chanAlloc = struct {
	mask uint32
	mx   sync.Mutex
}{mask: 0xffff_ffff}

// AllocChannel allocates a free channel in the controller. It returns an
// invalid channel if there is no free channel to be allocated. Use Channel.Free
// to free an unused channel.
func (d *Controller) AllocChannel() (ch Channel) {
	chanAlloc.mx.Lock()
	if chanAlloc.mask != 0 {
		mask := uint32(1)
		if chanAlloc.mask+1 == 0 {
			// Setup DMA before first use.
			RESETS := resets.RESETS()
			RESETS.RESET.ClearBits(resets.DMA) // remove reset
			for RESETS.RESET_DONE.LoadBits(resets.DMA) == 0 {
			}

			runtime.LockOSThread()
			pl, _ := rtos.SetPrivLevel(0)

			// Allow access from user mode
			for i := range d.seccfgCh {
				d.seccfgCh[i].Store(0b10)
			}
			for i := range d.seccfgIRQ {
				d.seccfgIRQ[i].Store(0b10)
			}
			d.seccfgMisc.Store(0b10_1010_1010)

			rtos.SetPrivLevel(pl)
			runtime.UnlockOSThread()

			chanAlloc.mask = 0xfffe
		} else {
			// Find a free channel.
			for chanAlloc.mask&mask == 0 {
				mask <<= 1
				ch.n++
			}
			chanAlloc.mask &^= mask
		}
		ch.d = d
	}
	chanAlloc.mx.Unlock()
	return
}

// Trig triggers the transfers on channels specified by the channels bitmask.
//
//go:nosplit
func (d *Controller) Trig(channels uint32) {
	d.multiChanTrig.Store(channels)
}

// RawIRQs returns the bitmask that represents the channels with an interrupt
// asserted internally in the DMA controller.
//
//go:nosplit
func (d *Controller) RawIRQs() uint32 {
	return d.irq[0].r.Load()
}

// ActiveIRQs returns the bitmask that represents the channels with an active
// interrupt requests to the CPU (only enabled interrupts are listed).
//
//go:nosplit
func (d *Controller) ActiveIRQs() uint32 {
	return d.irq[0].s.Load()
}

// ClearIRQs clears the interrupts for the channels specified by bitmask.
//
//go:nosplit
func (d *Controller) ClearIRQs(irqs uint32) {
	d.irq[0].r.Store(irqs)
}
