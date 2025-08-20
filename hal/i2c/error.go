// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c

import (
	"unsafe"

	"github.com/embeddedgo/device/bus/i2cbus"
)

// An Error wraps the value of TX_ABRT_SOURCE regisert.
type Error struct {
	Abort TX_ABRT_SOURCE
}

func (e *Error) Is(target error) bool {
	const ackFlags = ABRT_SBYTE_ACKDET | ABRT_HS_ACKDET | ABRT_GCALL_NOACK |
		ABRT_TXDATA_NOACK | ABRT_10ADDR2_NOACK | ABRT_10ADDR1_NOACK |
		ABRT_7B_ADDR_NOACK
	return target == i2cbus.ErrACK && e.Abort&ackFlags != 0
}

var abortStr = [...]string{
	0:  "7b_addr_noack",
	1:  "10addr1_noack",
	2:  "10addr2_noack",
	3:  "txdata_noack",
	4:  "gcall_noack",
	5:  "gcall_read",
	6:  "hs_ackdet",
	7:  "sbyte_ackdet",
	8:  "hs_norstrt",
	9:  "sbyte_norstrt",
	10: "10b_rd_norstrt",
	11: "master_dis",
	12: "lost",
	13: "slvflush_txfifo",
	14: "slv_arblost",
	15: "slvrd_intx",
	16: "user_abrt",
}

func (e *Error) Error() string {
	n := 0
	for i, s := range abortStr {
		if e.Abort>>uint(i)&1 != 0 {
			n += len(s) + 1
		}
	}
	prefix := "i2c: "
	buf := make([]byte, len(prefix), len(prefix)+n-1)
	copy(buf, prefix)
	for i, s := range abortStr {
		if e.Abort>>uint(i)&1 != 0 {
			buf = append(buf, s...)
			if len(buf) != cap(buf) {
				buf = append(buf, ',')
			}
		}
	}
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}
