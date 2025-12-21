// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"embedded/mmio"
	"structs"
)

const (
	imCap = 32
	numSM = 4
)

type Periph struct {
	_ structs.HostLayout

	CTRL              mmio.R32[CTRL]
	FSTAT             mmio.R32[FSTAT]
	FDEBUG            mmio.R32[FDEBUG]
	FLEVEL            mmio.R32[FLEVEL]
	TXF               [numSM]mmio.U32
	RXF               [numSM]mmio.U32
	_                 [2]uint32
	INPUT_SYNC_BYPASS mmio.U32
	DBG_PADOUT        mmio.U32
	DBG_PADOE         mmio.U32
	DBG_CFGINFO       mmio.R32[DBG_CFGINFO]
	INSTR_MEM         [imCap]mmio.R32[uint32]
	SM                [numSM]SMRegs
	RXF_PUTGET        [numSM][4]mmio.U32
	_                 [12]uint32
	GPIOBASE          mmio.U32
	INTR              mmio.R32[INTR]
	IRQ               [2]SIRQ
}

type SMRegs struct {
	_ structs.HostLayout

	CLKDIV    mmio.U32
	EXECCTRL  mmio.R32[EXECCTRL]
	SHIFTCTRL mmio.R32[SHIFTCTRL]
	ADDR      mmio.U32
	INSTR     mmio.U32
	PINCTRL   mmio.R32[PINCTRL]
}

type SIRQ struct {
	_ structs.HostLayout

	E mmio.R32[INTR]
	F mmio.R32[INTR]
	S mmio.R32[INTR]
}
