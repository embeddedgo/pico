// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "funcdata.h"
#include "textflag.h"

#define ICSR_ADDR 0xe000ed04
#define ICSR_PENDSVSET (1<<28)
#define NVIC_ISER0 0xE000E100
#define NVIC_IPR0 0xE000E400
#define RP2350_RAMEND (0x20000000 + 520*1024)
#define SIO_CPUID_ADDR 0xD0000000
#define SIO_DOORBELL_OUT_SET 0xD0000180
#define SIO_DOORBELL_IN_CLR 0xD000018C
#define SIO_FIFO_ADDR 0xD0000050
#define FIFO_ST 0
#define FIFO_WR 4
#define FIFO_RD 8
#define ST_VLD 1
#define ST_RDY 2

// _rt0_thumb_noos is the first function of the Embedded Go program
TEXT _rt0_thumb_noos(SB),NOSPLIT|NOFRAME,$0
	// Uncomment in case of truble to connect the debugger to the running program.
	//B  0(PC)  // gdb `set $pc += 4` to exit this loop

	// Disable exceptions until the CPU is ready for SVCall and PendSV.
	CPSID

	// Use the last 8 KB of RAM for two 4 KB stacks
	MOVW  $runtime·ramend(SB), R0
	MOVW  $RP2350_RAMEND, R13
	SUB   R0, R13, R0
	SRL   $1, R0  // stack size available for a single core
	CMP   $2048, R0
	BLT   0(PC)  // stuck here if the available stack size is less than 2 KiB

	// Set MSP and MSPLIMIT for this core
	MOVW  $SIO_CPUID_ADDR, R1
	MOVW  (R1), R1
	MUL   R1, R0, R2
	SUB   R2, R13
	SUB   R0, R13, R0
	MOVW  R0, MSPLIM

	// For now set PSPLIM to the beggining of free memory
	MOVW  $runtime·end(SB), R0
	MOVW  R0, PSPLIM

	// Cortex-M settings.
	MOVW       $0, R0                    // dummy RA
	MOVW       $runtime·vectors(SB), R1  // arg
	MOVM.DB.W  [R0-R1], (R13)
	BL         runtime·initCPU(SB)
	ADD        $8, R13

	// Enable SIO_IRQ_BELL (IRQ26) to support the runtime.preepmtOrWakeup
	MOVW  $(NVIC_IPR0+26), R0
	MOVW  $(14<<4), R1  // slightly higher than PendSV priority
	MOVB  R1, (R0)
	MOVW  $NVIC_ISER0, R0
	MOVW  $(1<<26), R1  // IRQ26
	MOVW  R1, (R0)

	// CPU execution paths diverge here.
	MOVW  $SIO_CPUID_ADDR, R1
	MOVW  (R1), R1
	CBZ   R1, cpu0

	// CPU1 waits until CPU0 prepares the tasker
	MOVW  $runtime·runOtherCPUs(SB), R0
waitTaskerReady:
	WFE
	DMB   MB_ISH
	MOVB  (R0), R1
	CMP   $0, R1
	BEQ   waitTaskerReady

	// CPU1 enters sheduler. Use PSP as SP (as threads do) to ensure the MSP
	// points to the top of the main stack (PSP will be overwritten anyway).
	BL    runtime·identcurcpu(SB)
	MOVW  R0, g
	MOVW  R13, R0
	MOVW  R0, PSP
	DSB
	MOVW  CONTROL, R0
	ORR   $2, R0  // use PSP as a stack pointer
	MOVW  R0, CONTROL
	ISB
	BL     runtime·curcpuSchedule(SB)  // raise the PendSV exception
	CPSIE  // enable exceptions (CPU1 should take PendSV immediately)
	B      0(PC)  // CPU1 never retuns here from PendSV

cpu0:
	// Clear memory and load the data segment from Flash.
	BL  runtime·initRAMfromROM(SB)
	DSB

	// Run CPU1
	MOVW       $0, R0                    // 0
	MOVW       $0, R1                    // 0
	MOVW       $1, R2                    // 1
	SLL        $29, R2, R3               // fake vectors_table
	MOVW       R3, R4                    // fake SP
	MOVW       $_rt0_thumb_noos(SB), R5  // entry
	MOVM.DB.W  [R0-R5], (R13)            // push commands to stack
	MOVW       $SIO_FIFO_ADDR, R2
cmdLoop:
	MOVW  (R13)(R0*4), R1  // load command from stack
	CBNZ  R1, waitTxReady
drainRx:
	MOVW  FIFO_ST(R2), R3
	AND   $ST_VLD, R3
	CBZ   R3, 3(PC)
	MOVW  FIFO_RD(R2), R3
	B     drainRx
	SEV
waitTxReady:
	MOVW   FIFO_ST(R2), R3
	AND.S  $ST_RDY, R3
	BEQ    waitTxReady
	MOVW   R1, FIFO_WR(R2)  // send command to CPU1
	DSB
	SEV  // wake up CPU1 that may sleep waiting for new data in FIFO
waitRxReady:
	WFE
	MOVW     FIFO_ST(R2), R3
	AND.S    $ST_VLD, R3
	BEQ      waitRxReady
	MOVW     FIFO_RD(R2), R3  // read response from CPU1
	CMP      R1, R3           // compare the response with the command sent
	MOVW.NE  $0, R0           // start from beginning
	ADD.EQ   $1, R0           // next command
	CMP      $6, R0
	BLT      cmdLoop
	ADD      $24, R13

	// Set the numer of available CPUs.
	MOVW  $2, R0
	MOVW  $runtime·ncpu(SB), R1
	MOVW  R0, (R1)

	// Go to the standard noos/thumb initialization code.
	MOVW  $runtime·ramend(SB), R0
	MOVW  $RP2350_RAMEND, R1
	B     runtime·rt0_go(SB)


// identcurcpu indetifies the current CPU and returns the pointer to its cpuctx
// in R0. It can clobber R0-R4,LR registers (other registers must be preserved).
TEXT runtime·identcurcpu(SB),NOSPLIT|NOFRAME,$0-0
	MOVW  $runtime·thetasker(SB), R0
	MOVW  (R0), R0  // allcpu is the first field of the runtime.tasker struct
	MOVW  $SIO_CPUID_ADDR, R1
	MOVW  (R1), R1
	MOVW  (R0)(R1*4), R0  // R0 = thetasker.allcpu[cpuid]
	RET


// func preepmtOrWakeup(cpuid int)
TEXT runtime·preemptOrWakeup(SB),NOSPLIT|NOFRAME,$0-4
	MOVW    cpuid+0(FP), R0
	CMP     $-1, R0
	RET.EQ  // wakeup the current CPU (no need to do anything)

	MOVW  $SIO_DOORBELL_OUT_SET, R0
	MOVW  $1, R1
	MOVW  R1, (R0)  // rise SIO_IRQ_BELL (IRQ26) on the opposite core
	RET


// SIO_IRQ_BELL handler
TEXT IRQ26_Handler(SB),NOSPLIT|NOFRAME,$0-0
	MOVW  $SIO_DOORBELL_IN_CLR, R0
	MOVW  $1, R1
	MOVW  R1, (R0)  // clear this IRQ

	MOVW  $ICSR_ADDR, R0
	MOVW  $ICSR_PENDSVSET, R1
	MOVW  R1, (R0)  // rise PendSV

	RET
