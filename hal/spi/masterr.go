// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import (
	"runtime"
	"unsafe"

	"github.com/embeddedgo/pico/hal/dma"
)

func (d *Master) SetRepW(repw uint16) {
	d.repw = repw
}

// Read implements io.Reader.
func (d *Master) Read(s []byte) (n int, err error) {
	if len(s) == 0 {
		return
	}
	pr := unsafe.Pointer(unsafe.SliceData(s))
	n = len(s)

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.rdma.IsValid() {
		read[uint8](d, pr, n)
	} else {
		pw := unsafe.Pointer(&d.repw)
		writeReadDMA(d, pw, pr, n, dma.S8b, dma.S8b|dma.IncW)
	}
	return
}

// Read16
func (d *Master) Read16(s []uint16) (n int, err error) {
	if len(s) == 0 {
		return
	}
	pr := unsafe.Pointer(unsafe.SliceData(s))
	n = len(s)

	// Use DMA only for long transfers. Short ones are handled by CPU.
	if n < minDMA || !d.rdma.IsValid() {
		read[uint16](d, pr, n)
	} else {
		pw := unsafe.Pointer(&d.repw)
		writeReadDMA(d, pw, pr, n, dma.S16b, dma.S16b|dma.IncW)
	}
	return
}

func read[T dataWord](d *Master, pr unsafe.Pointer, n int) {
	sz := int(unsafe.Sizeof(T(0)))
	p, slow, repw := d.p, d.slow, uint32(d.repw)
	nf := min(n, fifoLen)

	if d.wonly {
		drainRxFIFO(d)
	}

	// Fill the Tx FIFO.
	for range nf {
		p.DR.Store(repw)
	}
	n -= nf

	// Read and write.
	for end := unsafe.Add(pr, n*sz); pr != end; pr = unsafe.Add(pr, sz) {
		for p.SR.LoadBits(RNE) == 0 {
			if slow {
				runtime.Gosched()
			}
		}
		*(*T)(pr) = T(p.DR.Load())
		p.DR.Store(repw)
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

func (d *Master) ReadByte() (b byte, err error) {
	return byte(writeReadWord(d, uint32(d.repw))), nil
}

func (d *Master) ReadWord16() (w uint16, err error) {
	return uint16(writeReadWord(d, uint32(d.repw))), nil
}
