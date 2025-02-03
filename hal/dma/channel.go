// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dma

import "sync/atomic"

type Channel struct {
	d *Controller
	n int
}

func (c Channel) IsValid() bool {
	return c.d != nil
}

func (c Channel) Num() int {
	return c.n
}

// Free frees the channel so the Controller.AllocChannel can allocate it next
// time.
func (c Channel) Free() {
	mask := uint32(1) << uint(c.n)
	for {
		chs := atomic.LoadUint32(&chanMask)
		if atomic.CompareAndSwapUint32(&chanMask, chs, chs|mask) {
			break
		}
	}
}
