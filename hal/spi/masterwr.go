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
		writeRead8(d, pw, pr, n)
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
		writeRead16(d, pw, pr, n)
	} else {
		writeReadDMA(d, pw, pr, n, dma.S16b)
	}
	return
}

func writeRead8(d *Master, pw, pr unsafe.Pointer, n int) {
	p, slow := d.p, d.slow
	nf := min(n, fifoLen)

	if d.rdirty {
		drainRxFIFO(p)
		d.rdirty = false
	}

	// Fill the Tx FIFO.
	for end := unsafe.Add(pw, nf); pw != end; pw = unsafe.Add(pw, 1) {
		p.DR.Store(uint32(*(*uint8)(pw)))
	}
	n -= nf

	// Read and write.
	for end := unsafe.Add(pw, n); pw != end; pw = unsafe.Add(pw, 1) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*uint8)(pr) = uint8(p.DR.Load())
		p.DR.Store(uint32(*(*uint8)(pw)))
		pr = unsafe.Add(pr, 1)
	}

	// Read the remaining data
	for end := unsafe.Add(pr, nf); pr != end; pr = unsafe.Add(pr, 1) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*uint8)(pr) = uint8(p.DR.Load())
	}
}

func writeRead16(d *Master, pw, pr unsafe.Pointer, n int) {
	p, slow := d.p, d.slow
	nf := min(n, fifoLen)

	if d.rdirty {
		drainRxFIFO(p)
		d.rdirty = false
	}

	// Fill the Tx FIFO.
	for end := unsafe.Add(pw, nf*2); pw != end; pw = unsafe.Add(pw, 2) {
		p.DR.Store(uint32(*(*uint16)(pw)))
	}
	n -= nf

	// Read and write.
	for end := unsafe.Add(pw, n*2); pw != end; pw = unsafe.Add(pw, 2) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*uint16)(pr) = uint16(p.DR.Load())
		p.DR.Store(uint32(*(*uint16)(pw)))
		pr = unsafe.Add(pr, 2)
	}

	// Read the remaining data
	for end := unsafe.Add(pr, nf*2); pr != end; pr = unsafe.Add(pr, 2) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*uint16)(pr) = uint16(p.DR.Load())
	}
}

func writeReadDMA(d *Master, pw, pr unsafe.Pointer, n int, dmacfg dma.Config) {
	// TODO
}
