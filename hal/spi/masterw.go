// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import (
	"runtime"
	"unsafe"

	"github.com/embeddedgo/pico/hal/dma"
)

// Write implements io.Writer.
func (d *Master) Write(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	pw := unsafe.Pointer(unsafe.SliceData(s))
	n = len(s)

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.wdma.IsValid() {
		write[uint8](d, pw, n)
	} else {
		writeDMA(d, pw, n, dma.S8b|dma.IncR)
	}
	return
}

// WriteString implements io.StringWriter.
func (d *Master) WriteString(s string) (n int, err error) {
	return d.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// Write16
func (d *Master) Write16(s []uint16) (n int, err error) {
	if len(s) == 0 {
		return
	}
	pw := unsafe.Pointer(unsafe.SliceData(s))
	n = len(s)

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.wdma.IsValid() {
		write[uint16](d, pw, n)
	} else {
		writeDMA(d, pw, n, dma.S16b|dma.IncR)
	}
	return
}

func write[T dataWord](d *Master, pw unsafe.Pointer, n int) {
	sz := int(unsafe.Sizeof(T(0)))
	p, slow := d.p, d.slow
	d.wonly = true

	// Fill the FIFO fast if empty.
	if p.SR.LoadBits(TFE) != 0 {
		nf := min(n, fifoLen)
		for end := unsafe.Add(pw, nf*sz); pw != end; pw = unsafe.Add(pw, sz) {
			p.DR.Store(uint32(*(*T)(pw)))
		}
		n -= nf
	}

	// Slower "check before any write" way.
	for end := unsafe.Add(pw, n*sz); pw != end; pw = unsafe.Add(pw, sz) {
		for p.SR.LoadBits(TNF) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		p.DR.Store(uint32(*(*T)(pw)))
	}
}

func (d *Master) WriteByteN(b byte, n int) {
	if n <= 0 {
		return
	}

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.wdma.IsValid() {
		writeWordN(d, uint32(b), n)
	} else {
		writeDMA(d, unsafe.Pointer(&b), n, dma.S8b)
	}
}

func (d *Master) WriteWord16N(w uint16, n int) {
	if n <= 0 {
		return
	}

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.wdma.IsValid() {
		writeWordN(d, uint32(w), n)
	} else {
		writeDMA(d, unsafe.Pointer(&w), n, dma.S16b)
	}
}

func writeWordN(d *Master, w uint32, n int) {
	p, slow := d.p, d.slow
	d.wonly = true

	// Fill the FIFO fast if empty.
	if p.SR.LoadBits(TFE) != 0 {
		nf := min(n, fifoLen)
		for range nf {
			p.DR.Store(w)
		}
		n -= nf
	}

	// Slower "check before any write" way.
	for range n {
		for p.SR.LoadBits(TNF) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		p.DR.Store(w)
	}
}

func (d *Master) WriteByte(b byte) error {
	writeWord(d, uint32(b))
	return nil
}

func (d *Master) WriteWord16(w uint16) error {
	writeWord(d, uint32(w))
	return nil
}

func writeWord(d *Master, w uint32) {
	p, slow := d.p, d.slow
	d.wonly = true

	for p.SR.LoadBits(TNF) == 0 {
		if slow {
			runtime.Gosched()
		}
	}
	p.DR.Store(w)
}

func writeDMA(d *Master, pw unsafe.Pointer, n int, dmacfg dma.Config) {
	_writeDMA(d, uintptr(pw), n, dmacfg)
}

//go:uintptrescapes
func _writeDMA(d *Master, pw uintptr, n int, dmacfg dma.Config) {
	d.wonly = true
	d.done.Clear() // memory barrier
	wdma := d.wdma
	wdma.ClearIRQ()
	wdma.SetReadAddr(unsafe.Pointer(pw))
	wdma.SetTransCount(n, dma.Normal)
	wdma.SetConfigTrig(d.wdc|dmacfg, wdma)
	wdma.EnableIRQ(d.irqn)
	d.done.Sleep(-1)
}
