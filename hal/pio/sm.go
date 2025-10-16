// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"embedded/mmio"
	"structs"
	"unsafe"
)

type SM struct {
	_ structs.HostLayout

	CLKDIV    mmio.U32
	EXECCTRL  mmio.R32[EXECCTRL]
	SHIFTCTRL mmio.R32[SHIFTCTRL]
	ADDR      mmio.U32
	INSTR     mmio.U32
	PINCTRL   mmio.R32[PINCTRL]
}

func (sm *SM) PIO() *PIO {
	addr := uintptr(unsafe.Pointer(sm)) &^ (pioStep - 1)
	return (*PIO)(unsafe.Pointer(addr))
}
