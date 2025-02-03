// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dma

import (
	"embedded/mmio"
	"sync/atomic"
	"unsafe"
)

type Channel struct {
	d *Controller
	n int
}

func (c Channel) IsValid() bool {
	return c.d != nil
}

func (c Channel) Num() int {
	return c.n
}

// Free frees the channel so the Controller.AllocChannel can allocate it next
// time.
func (c Channel) Free() {
	mask := uint32(1) << uint(c.n)
	for {
		chs := atomic.LoadUint32(&chanMask)
		if atomic.CompareAndSwapUint32(&chanMask, chs, chs|mask) {
			break
		}
	}
}

func addr(r *mmio.U32) uintptr {
	return uintptr(r.Load())
}

func setAddr(r *mmio.U32, a unsafe.Pointer) {
	r.Store(uint32(uintptr(a)))
}

func transCount(r *mmio.U32) (cnt int, mode int8) {
	v := r.Load()
	return int(v & 0x0fff_ffff), int8(v >> 28)
}

func setTransCount(r *mmio.U32, cnt int, mode int8) {
	r.Store(uint32(cnt)&0x0fff_ffff | uint32(mode)<<28)
}

func (c Channel) ReadAddr() uintptr {
	return addr(&c.d.ch[c.n].readAddr)
}

func (c Channel) SetReadAddr(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].readAddr, a)
}

func (c Channel) SetReadAddrTrig(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].readAddrTrig3, a)
}

func (c Channel) WriteAddr() uintptr {
	return addr(&c.d.ch[c.n].writeAddr)
}

func (c Channel) SetWriteAddr(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].writeAddr, a)
}

func (c Channel) SetWriteAddrTrig(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].writeAddrTrig2, a)
}

// Transfer count mode
const (
	Normal      int8 = 0  // normal mode, compatible with the RP2040
	TriggerSelf int8 = 1  // re-triggers itself at the end of transfer
	Endless     int8 = 15 // perform an endless sequence of transfers
)

func (c Channel) TransCount() (cnt int, mode int8) {
	return transCount(&c.d.ch[c.n].transCount)
}

func (c Channel) SetTransCount(cnt int, mode int8) {
	setTransCount(&c.d.ch[c.n].transCount, cnt, mode)
}

func (c Channel) SetTransCountTrig(cnt int, mode int8) {
	setTransCount(&c.d.ch[c.n].transCountTrig1, cnt, mode)
}

type Ctrl uint32

// Ctlr bits
const (
	En       Ctrl = 1 << 0 // DMA channel enable
	HighPrio Ctrl = 1 << 1 // Referential treatment in issue scheduling

	DataSize Ctrl = 3 << 2 // Size of each bus transfer:
	Size8b   Ctrl = 0 << 2 // - byte
	Size16b  Ctrl = 1 << 2 // - half word
	Size32b  Ctrl = 2 << 2 // - word

	IncRead  Ctrl = 1 << 4 // Increment read address with each transfer
	ReadRev  Ctrl = 1 << 5 // Decremented read address rather than incremented
	IncWrite Ctrl = 1 << 6 // Increment write address with each transfer
	WriteRev Ctrl = 1 << 7 // Decremented write address rather than incremented

)

func RingSize(log2size int) (c Ctrl) {
	return Ctrl(log2size&15) << 8
}

func (c Ctrl) RingSize() (log2size int) {
	return int(c >> 8 & 15)
}
