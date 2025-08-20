// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c

import (
	"embedded/mmio"
	"structs"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/mmap"
	"github.com/embeddedgo/pico/p/resets"
)

type Periph struct {
	_ structs.HostLayout

	CON                mmio.R32[CON]
	TAR                mmio.R32[TAR]
	SAR                mmio.U32
	_                  uint32
	DATA_CMD           mmio.U32
	SS_SCL_HCNT        mmio.U32
	SS_SCL_LCNT        mmio.U32
	FS_SCL_HCNT        mmio.U32
	FS_SCL_LCNT        mmio.U32
	_                  [2]uint32
	INTR_STAT          mmio.R32[INTR]
	INTR_MASK          mmio.R32[INTR]
	RAW_INTR_STAT      mmio.R32[INTR]
	RX_TL              mmio.U32
	TX_TL              mmio.U32
	CLR_INTR           mmio.U32
	CLR_RX_UNDER       mmio.U32
	CLR_RX_OVER        mmio.U32
	CLR_TX_OVER        mmio.U32
	CLR_RD_REQ         mmio.U32
	CLR_TX_ABRT        mmio.U32
	CLR_RX_DONE        mmio.U32
	CLR_ACTIVITY       mmio.U32
	CLR_STOP_DET       mmio.U32
	CLR_START_DET      mmio.U32
	CLR_GEN_CALL       mmio.U32
	ENABLE             mmio.R32[ENABLE]
	STATUS             mmio.R32[STATUS]
	TXFLR              mmio.U32
	RXFLR              mmio.U32
	SDA_HOLD           mmio.R32[SDA_HOLD]
	TX_ABRT_SOURCE     mmio.R32[TX_ABRT_SOURCE]
	SLV_DATA_NACK_ONLY mmio.R32[SLV_DATA_NACK_ONLY]
	DMA_CR             mmio.R32[DMA_CR]
	DMA_TDLR           mmio.U32
	DMA_RDLR           mmio.U32
	SDA_SETUP          mmio.U32
	ACK_GENERAL_CALL   mmio.R32[ACK_GENERAL_CALL]
	ENABLE_STATUS      mmio.R32[ENABLE_STATUS]
	FS_SPKLEN          mmio.U32
	_                  uint32
	CLR_RESTART_DET    mmio.U32
	_                  [18]uint32
	COMP_PARAM_1       mmio.R32[COMP_PARAM_1]
	COMP_VERSION       mmio.U32
	COMP_TYPE          mmio.U32
}

// I2C returns the n-th instance of the I2C peripheral.
func I2C(n int) *Periph {
	if uint(n) > 1 {
		panic("wrong I2C number")
	}
	const base = mmap.I2C0_BASE
	const step = mmap.I2C1_BASE - mmap.I2C0_BASE
	return (*Periph)(unsafe.Pointer(base + uintptr(n)*step))
}

func num(p *Periph) int {
	const step = mmap.I2C1_BASE - mmap.I2C0_BASE
	return int((uintptr(unsafe.Pointer(p)) - mmap.I2C0_BASE) / step)
}

// SetReset allows to assert/deassert the reset signal to the I2C peripheral.
func (p *Periph) SetReset(assert bool) {
	internal.SetReset(resets.I2C0<<uint(num(p)), assert)
}
