// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

import (
	"embedded/rtos"
	"runtime"
	"sync"
	"unsafe"

	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/hal/system/clock"
)

// A Master is a driver to the SPI peripheral used in master mode.
type Master struct {
	sync.Mutex // helps in case of concurent use of Master (not used internally)

	p     *Periph
	slow  bool
	wonly bool
	repw  uint16

	rdma, wdma dma.Channel
	rdc, wdc   dma.Config

	irqn int
	done rtos.Note
}

// NewMaster returns a new master-mode driver for p. If valid DMA channels are
// given, the DMA will be used for bigger data transfers.
func NewMaster(p *Periph, rdma, wdma dma.Channel) *Master {
	irqn := int(system.NextCPU())
	d := &Master{p: p, rdma: rdma, wdma: wdma, irqn: irqn}
	reqAdd := dma.Config(num(p)) * (dma.SPI1_TX - dma.SPI0_TX)
	if rdma.IsValid() {
		d.rdc = dma.En | (dma.SPI0_RX + reqAdd)
		rdma.SetReadAddr(unsafe.Pointer(&p.DR))
	}
	if wdma.IsValid() {
		d.wdc = dma.En | (dma.SPI0_TX + reqAdd)
		wdma.SetWriteAddr(unsafe.Pointer(&p.DR))
	}
	return d
}

// Periph returns the underlying SPI peripheral.
func (d *Master) Periph() *Periph {
	return d.p
}

func (d *Master) Enable() {
	d.p.CR1.SetBits(SSE)
}

// Disable waits for the last bit of the last transfer to be sent and next it
// disables the SPI peripheral.
func (d *Master) Disable() {
	d.WaitTxDone()
	d.p.CR1.ClearBits(SSE)
}

func (d *Master) WaitTxDone() {
	p, slow := d.p, d.slow
	for p.SR.LoadBits(TFE|BSY) != TFE {
		if slow {
			runtime.Gosched()
		}
	}
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
	d.WaitTxDone()
	p := d.p
	p.CR0.StoreBits(FRF|SPH|SPO|DSS, CR0(cfg)) // only MS requires disabled SSP
	p.DMACR.Store(RXDMAE | TXDMAE)
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
	d.WaitTxDone()
	p := d.p
	p.CPSR.Store(uint32(cpsr))
	p.CR0.StoreBits(SCR, CR0(scr-1)<<SCRn) // only MS requires disabled SSP
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

// DMAISR should be configured as a DMA interrupt handler if DMA is used.
//
//go:nosplit
func (d *Master) DMAISR() {
	ch := d.rdma
	if d.wonly {
		ch = d.wdma
	}
	ch.DisableIRQ(d.irqn)
	d.done.Wakeup()
}

const (
	minDMA  = 32
	fifoLen = 8
)

func drainRxFIFO(d *Master) {
	p, slow := d.p, d.slow
	for {
		for p.SR.LoadBits(RNE) != 0 {
			p.DR.Load()
		}
		if p.SR.LoadBits(BSY) == 0 {
			break
		}
		if slow {
			runtime.Gosched()
		}
	}
	d.wonly = false
}

type dataWord interface{ ~uint8 | ~uint16 }
