// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import (
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/system/clock"
)

// A Master is a driver to the SPI peripheral used in master mode.
type Master struct {
	p      *Periph
	slow   bool
	rdirty bool
	repw   uint16

	rdma dma.Channel
	wdma dma.Channel
}

// NewMaster returns a new master-mode driver for p. If valid DMA channels are
// given, the DMA will be used for bigger data transfers.
func NewMaster(p *Periph, rxdma, txdma dma.Channel) *Master {
	return &Master{p: p, rdma: rxdma, wdma: txdma}
}

// Periph returns the underlying SPI peripheral.
func (d *Master) Periph() *Periph {
	return d.p
}

func (d *Master) Enable() {
	d.p.CR1.SetBits(SSE)
}

func (d *Master) Disable() {
	d.p.CR1.ClearBits(SSE)
}

type Config uint32

const (
	FrameFormat = Config(FRF) // Frame format:
	MSPI        = Config(0)   // - Motorola SPI
	SyncSerial  = Config(FTI) // - Texas Instruments Synchronous Serial
	Microwire   = Config(FNM) // - National Semiconductor Microwire

	// For Motorola SPI frame format:
	CPHA0 = Config(0)   // sample on leading edge.
	CPHA1 = Config(SPH) // sample on trailing edge.
	CPOL0 = Config(0)   // clock idle state is 0.
	CPOL1 = Config(SPO) // clock idle state is 1.

	WordLen = Config(15) // Data word length:
	Word4b  = Config(3)  // - 4 bit
	Word5b  = Config(4)  // - 5 bit
	Word6b  = Config(5)  // - 6 bit
	Word7b  = Config(6)  // - 7 bit
	Word8b  = Config(7)  // - 8 bit
	Word9b  = Config(8)  // - 9 bit
	Word10b = Config(9)  // - 10 bit
	Word11b = Config(10) // - 11 bit
	Word12b = Config(11) // - 12 bit
	Word13b = Config(12) // - 13 bit
	Word14b = Config(13) // - 14 bit
	Word15b = Config(14) // - 15 bit
	Word16b = Config(15) // - 16 bit
)

func (d *Master) Config() Config {
	return Config(d.p.CR0.LoadBits(FRF | SPH | SPO | DSS))
}

func (d *Master) SetConfig(cfg Config) {
	p := d.p
	cr1 := p.CR1.Load()
	p.CR1.Store(cr1 &^ SSE) // disable SPI
	p.CR0.StoreBits(FRF|SPH|SPO|DSS, CR0(cfg))
	p.DMACR.Store(3)
	p.CR1.Store(cr1)
}

func (d *Master) Baudrate() int {
	p := d.p
	scr := uint(p.CR0.LoadBits(SCR) >> SCRn)
	cpsr := uint(p.CPSR.LoadBits(0xff))
	div := int64((scr + 1) * cpsr)
	return int((uint(2*clock.PERI.Freq()/div) + 1) / 2)
}

func (d *Master) SetBaudrate(baudrate int) (actual int) {
	if baudrate <= 0 {
		return -1
	}
	periHz := clock.PERI.Freq()
	div := uint((periHz + int64(baudrate-1)) / int64(baudrate))
	if div < 2 {
		return -1
	}
	cpsr := uint(2)
	var scr uint
	for {
		scr = (div + cpsr - 1) / cpsr
		if scr <= 256 {
			// We found the smalest cpsr and largest scr so the cpsr*scr is
			// close to the div (cpsr*scr >= div). It may not be the closest
			// possible combination but we use it as is.
			break
		}
		if cpsr += 2; cpsr > 254 {
			return -1
		}
	}
	p := d.p
	cr1 := p.CR1.Load()
	p.CR1.Store(cr1 &^ SSE) // disable SPI
	p.CPSR.Store(uint32(cpsr))
	p.CR0.StoreBits(SCR, CR0(scr-1)<<SCRn)
	p.CR1.Store(cr1)
	div = scr * cpsr
	actual = int((uint(2*periHz/int64(div)) + 1) / 2)
	d.slow = actual <= 1e5
	return
}

// Setup resets the underlying SPI peripheral and configures it according to
// the master driver needs. Next it calls the SetConfig and SetBaudrate methods
// with the provided arguments and enables the peripheral.
func (d *Master) Setup(cfg Config, baudrate int) (actualBaud int) {
	p := d.p
	p.SetReset(true)
	p.SetReset(false)
	d.SetConfig(cfg)
	actualBaud = d.SetBaudrate(baudrate)
	d.Enable()
	return
}

const (
	minDMA  = 32
	fifoLen = 8
)

func drainRxFIFO(p *Periph) {
	for {
		for p.SR.LoadBits(RNE) != 0 {
			p.DR.Load()
		}
		if p.SR.LoadBits(BSY) == 0 {
			return
		}
	}
}

type dataWord interface{ ~uint8 | ~uint16 }
