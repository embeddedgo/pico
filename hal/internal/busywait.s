// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// func BusyWaitAtLeastCycles(uint n)
TEXT ·BusyWaitAtLeastCycles(SB),NOSPLIT|NOFRAME,$0-4
	MOVW  n+0(FP), R0
	SUB.S $3, R0
	BCS -1(PC)
	RET
