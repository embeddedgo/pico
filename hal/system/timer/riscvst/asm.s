// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

#define ICSR_ADDR 0xe000ed04
#define ICSR_PENDSVSET (1<<28)

#define MTIMECMPH_ADDR (0xd0000000 + 0x1bc)

// SIO_IRQ_MTIMECMP handler
TEXT IRQ29_Handler(SB),NOSPLIT|NOFRAME,$0-0
	// Set PendSV bit first to avoid DSB but ensure exception tail-chaining
	MOVW  $ICSR_ADDR, R0
	MOVW  $ICSR_PENDSVSET, R1
	MOVW  R1, (R0)
	//SEV   // see ARM Errata 563915

	// Clear this IRQ.
	MOVW  $MTIMECMPH_ADDR, R0
	MOVW  $-1, R1
	MOVW  R1, (R0)

	RET
