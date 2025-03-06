// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package system

import (
	"embedded/rtos"
	"sync/atomic"
)

const ncpus = 2 // NextCPU requires ncpus to be a power of 2

var cpu uint32

// NextCPU returns the interrupt context suitable for the rtos.IRQ.Enable
// function. Use it to evenly (in terms of quantity, not workload) distribute
// interrupts between available CPUs / cores.
func NextCPU() rtos.IntCtx {
	return rtos.IntCtx(atomic.AddUint32(&cpu, 1) & (ncpus - 1))
}
