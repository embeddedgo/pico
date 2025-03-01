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
		write8(d, pw, n)
	} else {
		writeDMA(d, pw, n, dma.S8b)
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
		write16(d, pw, n)
	} else {
		writeDMA(d, pw, n, dma.S16b)
	}
	return
}

func write8(d *Master, pw unsafe.Pointer, n int) {
	p, slow := d.p, d.slow

	// Fill the FIFO fast if empty.
	if p.SR.LoadBits(TFE) != 0 {
		nf := min(n, fifoLen)
		for end := unsafe.Add(pw, nf); pw != end; pw = unsafe.Add(pw, 1) {
			p.DR.Store(uint32(*(*uint8)(pw)))
		}
		n -= nf
	}

	// Slower "check before any write" way.
	for end := unsafe.Add(pw, n); pw != end; pw = unsafe.Add(pw, 1) {
		for p.SR.LoadBits(TNF) != 0 {
			if slow {
				runtime.Gosched()
			}
		}
		p.DR.Store(uint32(*(*uint8)(pw)))
	}

	d.rdirty = true // we left unread garbage
}

func write16(d *Master, pw unsafe.Pointer, n int) {
	p, slow := d.p, d.slow

	// Fill the FIFO fast if empty.
	if p.SR.LoadBits(TFE) != 0 {
		nf := min(n, fifoLen)
		for end := unsafe.Add(pw, nf*2); pw != end; pw = unsafe.Add(pw, 2) {
			p.DR.Store(uint32(*(*uint16)(pw)))
		}
		n -= nf
	}

	// Slower "check before any write" way.
	for end := unsafe.Add(pw, n*2); pw != end; pw = unsafe.Add(pw, 2) {
		for p.SR.LoadBits(TNF) != 0 {
			if slow {
				runtime.Gosched()
			}
		}
		p.DR.Store(uint32(*(*uint16)(pw)))
	}

	d.rdirty = true // we left unread garbage
}

func writeDMA(d *Master, pw unsafe.Pointer, n int, dmacfg dma.Config) {
	// TODO
}
