// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dma

import (
	"embedded/mmio"
	"runtime"
	"sync/atomic"
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
	RESETS := resets.RESETS()
	RESETS.RESET.ClearBits(resets.DMA)
	for RESETS.RESET_DONE.LoadBits(resets.DMA) == 0 {
		runtime.Gosched()
	}
}

func DMA(n int) *Controller {
	if n != 0 {
		panic("wrong DMA number")
	}
	return (*Controller)(unsafe.Pointer(mmap.DMA_BASE))
}

func (d *Controller) Channel(n int) Channel {
	return Channel{d, n}
}

var chanMask uint32 = 0xffff

// AllocChannel allocates a free channel in the controller. It returns an
// invalid channel if there is no free channel to be allocated. Use Channel.Free
// to free an unused channel.
func (d *Controller) AllocChannel() Channel {
	for {
		chs := atomic.LoadUint32(&chanMask)
		if chs == 0 {
			return Channel{}
		}
		n := 15
		mask := uint32(1) << uint(n)
		for chs&mask == 0 {
			mask >>= 1
			n--
		}
		if atomic.CompareAndSwapUint32(&chanMask, chs, chs&^mask) {
			return Channel{d, n}
		}
	}
}
