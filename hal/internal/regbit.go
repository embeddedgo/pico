// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"embedded/mmio"
	"unsafe"
)

// Atomic set, clear, xor bit operations using alias adresses. Not supported by
// SIO, OTP, CM33 PPB except Pi-specific registers on the EPPB.

func AtomicSet[T mmio.T32](r *mmio.R32[T], mask T) {
	(*mmio.R32[T])(unsafe.Pointer(r.Addr() + 0x2000)).Store(mask)
}

func AtomicClear[T mmio.T32](r *mmio.R32[T], mask T) {
	(*mmio.R32[T])(unsafe.Pointer(r.Addr() + 0x3000)).Store(mask)
}

func AtomicXor[T mmio.T32](r *mmio.R32[T], mask T) {
	(*mmio.R32[T])(unsafe.Pointer(r.Addr() + 0x1000)).Store(mask)
}

func AtomicMod[T mmio.T32](r *mmio.R32[T], mask, prev, new T) {
	(*mmio.R32[T])(unsafe.Pointer(r.Addr() + 0x1000)).Store((prev ^ new) & mask)
}

func AtomicSetU32(r *mmio.U32, mask uint32) {
	(*mmio.U32)(unsafe.Pointer(r.Addr() + 0x2000)).Store(mask)
}

func AtomicClearU32(r *mmio.U32, mask uint32) {
	(*mmio.U32)(unsafe.Pointer(r.Addr() + 0x3000)).Store(mask)
}

func AtomicXorU32(r *mmio.U32, mask uint32) {
	(*mmio.U32)(unsafe.Pointer(r.Addr() + 0x1000)).Store(mask)
}

func AtomicModU32(r *mmio.U32, mask, prev, new uint32) {
	(*mmio.U32)(unsafe.Pointer(r.Addr() + 0x1000)).Store((prev ^ new) & mask)
}
