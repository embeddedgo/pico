// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iomux

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/p/mmap"
)

type padsbank struct {
	voltage mmio.R32[uint32]
	pad     [50]mmio.R32[uint32]
}

func pb() *padsbank {
	return (*padsbank)(unsafe.Pointer(mmap.PADS_BANK0_BASE))
}

type iobank struct {
	gpio                              [48]sgpio
	_                                 [32]uint32
	irqsummary_proc0_secure           [2]mmio.R32[uint32]
	irqsummary_proc0_nonsecure        [2]mmio.R32[uint32]
	irqsummary_proc1_secure           [2]mmio.R32[uint32]
	irqsummary_proc1_nonsecure        [2]mmio.R32[uint32]
	irqsummary_dormant_wake_secure    [2]mmio.R32[uint32]
	irqsummary_dormant_wake_nonsecure [2]mmio.R32[uint32]
	intr                              [6]mmio.R32[uint32]
	proc0_inte                        [6]mmio.R32[uint32]
	proc0_intf                        [6]mmio.R32[uint32]
	proc0_ints                        [6]mmio.R32[uint32]
	proc1_inte                        [6]mmio.R32[uint32]
	proc1_intf                        [6]mmio.R32[uint32]
	proc1_ints                        [6]mmio.R32[uint32]
	dormant_wake_inte                 [6]mmio.R32[uint32]
	dormant_wake_intf                 [6]mmio.R32[uint32]
	dormant_wake_ints                 [6]mmio.R32[uint32]
}

type sgpio struct {
	status mmio.R32[uint32]
	ctrl   mmio.R32[uint32]
}

func ib() *iobank {
	return (*iobank)(unsafe.Pointer(uintptr(mmap.IO_BANK0_BASE)))
}
