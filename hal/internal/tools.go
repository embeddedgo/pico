// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"embedded/rtos"
	"sync/atomic"
)

func BusyWaitAtLeastCycles(n uint)

var cpu uint32

func NextCPU() rtos.IntCtx {
	return rtos.IntCtx(atomic.AddUint32(&cpu, 1))
}
