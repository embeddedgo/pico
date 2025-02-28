// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spi

type CR0 uint32

// CR0 bits
const (
	DSS  CR0 = 0x0F << 0 //+ Data Size Select:
	DS4  CR0 = 0x03 << 0 //  4-bit data
	DS5  CR0 = 0x04 << 0 //  5-bit data
	DS6  CR0 = 0x05 << 0 //  6-bit data
	DS7  CR0 = 0x06 << 0 //  7-bit data
	DS8  CR0 = 0x07 << 0 //  8-bit data
	DS9  CR0 = 0x08 << 0 //  9-bit data
	DS10 CR0 = 0x09 << 0 //  10-bit data
	DS11 CR0 = 0x0A << 0 //  11-bit data
	DS12 CR0 = 0x0B << 0 //  12-bit data
	DS13 CR0 = 0x0C << 0 //  13-bit data
	DS14 CR0 = 0x0D << 0 //  14-bit data
	DS15 CR0 = 0x0E << 0 //  15-bit data
	DS16 CR0 = 0x0F << 0 //  16-bit data
	FRF  CR0 = 0x03 << 4 //+ Frame format:
	FMO  CR0 = 0x00 << 4 //  Motorola SPI frame format
	FTI  CR0 = 0x01 << 4 //  TI synchronous serial frame format
	FNM  CR0 = 0x02 << 4 //  National Microwire frame format
	SPO  CR0 = 0x01 << 6 //+ SSPCLKOUT polarity for Motorola SPI frame format.
	SPH  CR0 = 0x01 << 7 //+ SSPCLKOUT phase for Motorola SPI frame format.
	SCR  CR0 = 0xFF << 8 //+ Serial clock rate.
)

const (
	DSSn = 0
	FRFn = 4
	SPOn = 6
	SPHn = 7
	SCRn = 8
)

type CR1 uint32

// CR1 bits
const (
	LBM CR1 = 0x01 << 0 //+ Loop back mode.
	SSE CR1 = 0x01 << 1 //+ Synchronous serial port enable.
	MS  CR1 = 0x01 << 2 //+ Slave mode.
	SOD CR1 = 0x01 << 3 //+ Slave-mode output disable.
)

type SR uint32

// SR bits
const (
	TFE SR = 0x01 << 0 //+ Transmit FIFO empty.
	TNF SR = 0x01 << 1 //+ Transmit FIFO not full.
	RNE SR = 0x01 << 2 //+ Receive FIFO not empty.
	RFF SR = 0x01 << 3 //+ Receive FIFO full.
	BSY SR = 0x01 << 4 //+ PrimeCell SSP busy flag.
)

type INT uint32

// INT bits
const (
	RORI INT = 0x01 << 0 //+ Receive overrun interrupt.
	RTI  INT = 0x01 << 1 //+ Receive timeout interrupt.
	RXI  INT = 0x01 << 2 //+ Receive FIFO interrupt.
	TXI  INT = 0x01 << 3 //+ Transmit FIFO interrupt.
)

type DMACR uint32

// DMACR bits
const (
	RXDMAE DMACR = 0x01 << 0 //+ Receive DMA Enable.
	TXDMAE DMACR = 0x01 << 1 //+ Transmit DMA Enable.
)

// PERIPHID1 bits
const (
	PARTNUMBER1 uint32 = 0x0F << 0 //+ These bits read back as 0x0
	DESIGNER0   uint32 = 0x0F << 4 //+ These bits read back as 0x1
)

const (
	PARTNUMBER1n = 0
	DESIGNER0n   = 4
)

// PERIPHID2 bits
const (
	DESIGNER1 uint32 = 0x0F << 0 //+ These bits read back as 0x4
	REVISION  uint32 = 0x0F << 4 //+ These bits return the peripheral revision
)

const (
	DESIGNER1n = 0
	REVISIONn  = 4
)
