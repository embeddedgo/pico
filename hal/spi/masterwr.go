// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import (
	"runtime"
	"unsafe"

	"github.com/embeddedgo/pico/hal/dma"
)

// WriteRead writes n = min(len(out), len(in)) bytes to the transmit FIFO. At
// the same time it reads n bytes into in.
func (d *Master) WriteRead(out, in []byte) (n int) {
	n = min(len(out), len(in))
	if n == 0 {
		return
	}
	pw := unsafe.Pointer(unsafe.SliceData(out))
	pr := unsafe.Pointer(unsafe.SliceData(in))

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.rdma.IsValid() || !d.wdma.IsValid() {
		writeRead[uint8](d, pw, pr, n)
	} else {
		writeReadDMA(d, pw, pr, n, dma.S8b)
	}
	return
}

// WriteStringRead works like WriteRead but writes a string.
func (d *Master) WriteStringRead(out string, in []byte) int {
	return d.WriteRead(unsafe.Slice(unsafe.StringData(out), len(out)), in)
}

// WriteRead16 writes n = min(len(out), len(in)) 16-bit words to the transmit
// FIFO. At the same time it reads n words into in.
func (d *Master) WriteRead16(out, in []uint16) (n int) {
	n = min(len(out), len(in))
	if n == 0 {
		return
	}
	pw := unsafe.Pointer(unsafe.SliceData(out))
	pr := unsafe.Pointer(unsafe.SliceData(in))

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.rdma.IsValid() || !d.wdma.IsValid() {
		writeRead[uint16](d, pw, pr, n)
	} else {
		writeReadDMA(d, pw, pr, n, dma.S16b)
	}
	return
}

func writeRead[T dataWord](d *Master, pw, pr unsafe.Pointer, n int) {
	sz := int(unsafe.Sizeof(T(0)))
	p, slow := d.p, d.slow
	nf := min(n, fifoLen)

	if d.rdirty {
		drainRxFIFO(p)
		d.rdirty = false
	}

	// Fill the Tx FIFO.
	for end := unsafe.Add(pw, nf*sz); pw != end; pw = unsafe.Add(pw, sz) {
		p.DR.Store(uint32(*(*T)(pw)))
	}
	n -= nf

	// Read and write.
	for end := unsafe.Add(pw, n*sz); pw != end; pw = unsafe.Add(pw, sz) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*T)(pr) = T(p.DR.Load())
		p.DR.Store(uint32(*(*T)(pw)))
		pr = unsafe.Add(pr, sz)
	}

	// Read the remaining data
	for end := unsafe.Add(pr, nf*sz); pr != end; pr = unsafe.Add(pr, sz) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*T)(pr) = T(p.DR.Load())
	}
}

func writeReadDMA(d *Master, pw, pr unsafe.Pointer, n int, dmacfg dma.Config) {
	// TODO
}

func (d *Master) WriteReadByte(b byte) byte {
	return byte(writeReadWord(d, uint32(b)))
}

func (d *Master) WriteReadWord16(w uint16) uint16 {
	return uint16(writeReadWord(d, uint32(w)))
}

func writeReadWord(d *Master, w uint32) uint32 {
	p, slow := d.p, d.slow

	if d.rdirty {
		drainRxFIFO(p)
		d.rdirty = false
	}

	p.DR.Store(w) // the Tx FIFO is empty
	for p.SR.LoadBits(RNE) == 0 {
		if slow {
			runtime.Gosched()
		}
	}
	return p.DR.Load()
}
