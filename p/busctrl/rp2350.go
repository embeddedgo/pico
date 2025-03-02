// DO NOT EDIT THIS FILE. GENERATED BY svdxgen.

//go:build rp2350

// Package busctrl provides access to the registers of the BUSCTRL peripheral.
//
// Instances:
//
//	BUSCTRL  BUSCTRL_BASE  -  -  Register block for busfabric control signals and performance counters
//
// Registers:
//
//	0x000 32  BUS_PRIORITY      Set the priority of each master for bus arbitration.
//	0x004 32  BUS_PRIORITY_ACK  Bus priority acknowledge
//	0x008 32  PERFCTR_EN        Enable the performance counters. If 0, the performance counters do not increment. This can be used to precisely start/stop event sampling around the profiled section of code. The performance counters are initially disabled, to save energy.
//	0x00C 32  PERFCTR0          Bus fabric performance counter 0
//	0x010 32  PERFSEL0          Bus fabric performance event select for PERFCTR0
//	0x014 32  PERFCTR1          Bus fabric performance counter 1
//	0x018 32  PERFSEL1          Bus fabric performance event select for PERFCTR1
//	0x01C 32  PERFCTR2          Bus fabric performance counter 2
//	0x020 32  PERFSEL2          Bus fabric performance event select for PERFCTR2
//	0x024 32  PERFCTR3          Bus fabric performance counter 3
//	0x028 32  PERFSEL3          Bus fabric performance event select for PERFCTR3
//
// Import:
//
//	github.com/embeddedgo/pico/p/mmap
package busctrl

const (
	PROC0 BUS_PRIORITY = 0x01 << 0  //+ 0 - low priority, 1 - high priority
	PROC1 BUS_PRIORITY = 0x01 << 4  //+ 0 - low priority, 1 - high priority
	DMA_R BUS_PRIORITY = 0x01 << 8  //+ 0 - low priority, 1 - high priority
	DMA_W BUS_PRIORITY = 0x01 << 12 //+ 0 - low priority, 1 - high priority
)

const (
	PROC0n = 0
	PROC1n = 4
	DMA_Rn = 8
	DMA_Wn = 12
)

const (
	PERFSEL0                    PERFSEL0 = 0x7F << 0 //+ Select an event for PERFCTR0. For each downstream port of the main crossbar, four events are available: ACCESS, an access took place; ACCESS_CONTESTED, an access took place that previously stalled due to contention from other masters; STALL_DOWNSTREAM, count cycles where any master stalled due to a stall on the downstream bus; STALL_UPSTREAM, count cycles where any master stalled for any reason, including contention from other masters.
	SIOB_PROC1_STALL_UPSTREAM   PERFSEL0 = 0x00 << 0
	SIOB_PROC1_STALL_DOWNSTREAM PERFSEL0 = 0x01 << 0
	SIOB_PROC1_ACCESS_CONTESTED PERFSEL0 = 0x02 << 0
	SIOB_PROC1_ACCESS           PERFSEL0 = 0x03 << 0
	SIOB_PROC0_STALL_UPSTREAM   PERFSEL0 = 0x04 << 0
	SIOB_PROC0_STALL_DOWNSTREAM PERFSEL0 = 0x05 << 0
	SIOB_PROC0_ACCESS_CONTESTED PERFSEL0 = 0x06 << 0
	SIOB_PROC0_ACCESS           PERFSEL0 = 0x07 << 0
	APB_STALL_UPSTREAM          PERFSEL0 = 0x08 << 0
	APB_STALL_DOWNSTREAM        PERFSEL0 = 0x09 << 0
	APB_ACCESS_CONTESTED        PERFSEL0 = 0x0A << 0
	APB_ACCESS                  PERFSEL0 = 0x0B << 0
	FASTPERI_STALL_UPSTREAM     PERFSEL0 = 0x0C << 0
	FASTPERI_STALL_DOWNSTREAM   PERFSEL0 = 0x0D << 0
	FASTPERI_ACCESS_CONTESTED   PERFSEL0 = 0x0E << 0
	FASTPERI_ACCESS             PERFSEL0 = 0x0F << 0
	SRAM9_STALL_UPSTREAM        PERFSEL0 = 0x10 << 0
	SRAM9_STALL_DOWNSTREAM      PERFSEL0 = 0x11 << 0
	SRAM9_ACCESS_CONTESTED      PERFSEL0 = 0x12 << 0
	SRAM9_ACCESS                PERFSEL0 = 0x13 << 0
	SRAM8_STALL_UPSTREAM        PERFSEL0 = 0x14 << 0
	SRAM8_STALL_DOWNSTREAM      PERFSEL0 = 0x15 << 0
	SRAM8_ACCESS_CONTESTED      PERFSEL0 = 0x16 << 0
	SRAM8_ACCESS                PERFSEL0 = 0x17 << 0
	SRAM7_STALL_UPSTREAM        PERFSEL0 = 0x18 << 0
	SRAM7_STALL_DOWNSTREAM      PERFSEL0 = 0x19 << 0
	SRAM7_ACCESS_CONTESTED      PERFSEL0 = 0x1A << 0
	SRAM7_ACCESS                PERFSEL0 = 0x1B << 0
	SRAM6_STALL_UPSTREAM        PERFSEL0 = 0x1C << 0
	SRAM6_STALL_DOWNSTREAM      PERFSEL0 = 0x1D << 0
	SRAM6_ACCESS_CONTESTED      PERFSEL0 = 0x1E << 0
	SRAM6_ACCESS                PERFSEL0 = 0x1F << 0
	SRAM5_STALL_UPSTREAM        PERFSEL0 = 0x20 << 0
	SRAM5_STALL_DOWNSTREAM      PERFSEL0 = 0x21 << 0
	SRAM5_ACCESS_CONTESTED      PERFSEL0 = 0x22 << 0
	SRAM5_ACCESS                PERFSEL0 = 0x23 << 0
	SRAM4_STALL_UPSTREAM        PERFSEL0 = 0x24 << 0
	SRAM4_STALL_DOWNSTREAM      PERFSEL0 = 0x25 << 0
	SRAM4_ACCESS_CONTESTED      PERFSEL0 = 0x26 << 0
	SRAM4_ACCESS                PERFSEL0 = 0x27 << 0
	SRAM3_STALL_UPSTREAM        PERFSEL0 = 0x28 << 0
	SRAM3_STALL_DOWNSTREAM      PERFSEL0 = 0x29 << 0
	SRAM3_ACCESS_CONTESTED      PERFSEL0 = 0x2A << 0
	SRAM3_ACCESS                PERFSEL0 = 0x2B << 0
	SRAM2_STALL_UPSTREAM        PERFSEL0 = 0x2C << 0
	SRAM2_STALL_DOWNSTREAM      PERFSEL0 = 0x2D << 0
	SRAM2_ACCESS_CONTESTED      PERFSEL0 = 0x2E << 0
	SRAM2_ACCESS                PERFSEL0 = 0x2F << 0
	SRAM1_STALL_UPSTREAM        PERFSEL0 = 0x30 << 0
	SRAM1_STALL_DOWNSTREAM      PERFSEL0 = 0x31 << 0
	SRAM1_ACCESS_CONTESTED      PERFSEL0 = 0x32 << 0
	SRAM1_ACCESS                PERFSEL0 = 0x33 << 0
	SRAM0_STALL_UPSTREAM        PERFSEL0 = 0x34 << 0
	SRAM0_STALL_DOWNSTREAM      PERFSEL0 = 0x35 << 0
	SRAM0_ACCESS_CONTESTED      PERFSEL0 = 0x36 << 0
	SRAM0_ACCESS                PERFSEL0 = 0x37 << 0
	XIP_MAIN1_STALL_UPSTREAM    PERFSEL0 = 0x38 << 0
	XIP_MAIN1_STALL_DOWNSTREAM  PERFSEL0 = 0x39 << 0
	XIP_MAIN1_ACCESS_CONTESTED  PERFSEL0 = 0x3A << 0
	XIP_MAIN1_ACCESS            PERFSEL0 = 0x3B << 0
	XIP_MAIN0_STALL_UPSTREAM    PERFSEL0 = 0x3C << 0
	XIP_MAIN0_STALL_DOWNSTREAM  PERFSEL0 = 0x3D << 0
	XIP_MAIN0_ACCESS_CONTESTED  PERFSEL0 = 0x3E << 0
	XIP_MAIN0_ACCESS            PERFSEL0 = 0x3F << 0
	ROM_STALL_UPSTREAM          PERFSEL0 = 0x40 << 0
	ROM_STALL_DOWNSTREAM        PERFSEL0 = 0x41 << 0
	ROM_ACCESS_CONTESTED        PERFSEL0 = 0x42 << 0
	ROM_ACCESS                  PERFSEL0 = 0x43 << 0
)

const (
	PERFSEL0n = 0
)

const (
	PERFSEL1                    PERFSEL1 = 0x7F << 0 //+ Select an event for PERFCTR1. For each downstream port of the main crossbar, four events are available: ACCESS, an access took place; ACCESS_CONTESTED, an access took place that previously stalled due to contention from other masters; STALL_DOWNSTREAM, count cycles where any master stalled due to a stall on the downstream bus; STALL_UPSTREAM, count cycles where any master stalled for any reason, including contention from other masters.
	SIOB_PROC1_STALL_UPSTREAM   PERFSEL1 = 0x00 << 0
	SIOB_PROC1_STALL_DOWNSTREAM PERFSEL1 = 0x01 << 0
	SIOB_PROC1_ACCESS_CONTESTED PERFSEL1 = 0x02 << 0
	SIOB_PROC1_ACCESS           PERFSEL1 = 0x03 << 0
	SIOB_PROC0_STALL_UPSTREAM   PERFSEL1 = 0x04 << 0
	SIOB_PROC0_STALL_DOWNSTREAM PERFSEL1 = 0x05 << 0
	SIOB_PROC0_ACCESS_CONTESTED PERFSEL1 = 0x06 << 0
	SIOB_PROC0_ACCESS           PERFSEL1 = 0x07 << 0
	APB_STALL_UPSTREAM          PERFSEL1 = 0x08 << 0
	APB_STALL_DOWNSTREAM        PERFSEL1 = 0x09 << 0
	APB_ACCESS_CONTESTED        PERFSEL1 = 0x0A << 0
	APB_ACCESS                  PERFSEL1 = 0x0B << 0
	FASTPERI_STALL_UPSTREAM     PERFSEL1 = 0x0C << 0
	FASTPERI_STALL_DOWNSTREAM   PERFSEL1 = 0x0D << 0
	FASTPERI_ACCESS_CONTESTED   PERFSEL1 = 0x0E << 0
	FASTPERI_ACCESS             PERFSEL1 = 0x0F << 0
	SRAM9_STALL_UPSTREAM        PERFSEL1 = 0x10 << 0
	SRAM9_STALL_DOWNSTREAM      PERFSEL1 = 0x11 << 0
	SRAM9_ACCESS_CONTESTED      PERFSEL1 = 0x12 << 0
	SRAM9_ACCESS                PERFSEL1 = 0x13 << 0
	SRAM8_STALL_UPSTREAM        PERFSEL1 = 0x14 << 0
	SRAM8_STALL_DOWNSTREAM      PERFSEL1 = 0x15 << 0
	SRAM8_ACCESS_CONTESTED      PERFSEL1 = 0x16 << 0
	SRAM8_ACCESS                PERFSEL1 = 0x17 << 0
	SRAM7_STALL_UPSTREAM        PERFSEL1 = 0x18 << 0
	SRAM7_STALL_DOWNSTREAM      PERFSEL1 = 0x19 << 0
	SRAM7_ACCESS_CONTESTED      PERFSEL1 = 0x1A << 0
	SRAM7_ACCESS                PERFSEL1 = 0x1B << 0
	SRAM6_STALL_UPSTREAM        PERFSEL1 = 0x1C << 0
	SRAM6_STALL_DOWNSTREAM      PERFSEL1 = 0x1D << 0
	SRAM6_ACCESS_CONTESTED      PERFSEL1 = 0x1E << 0
	SRAM6_ACCESS                PERFSEL1 = 0x1F << 0
	SRAM5_STALL_UPSTREAM        PERFSEL1 = 0x20 << 0
	SRAM5_STALL_DOWNSTREAM      PERFSEL1 = 0x21 << 0
	SRAM5_ACCESS_CONTESTED      PERFSEL1 = 0x22 << 0
	SRAM5_ACCESS                PERFSEL1 = 0x23 << 0
	SRAM4_STALL_UPSTREAM        PERFSEL1 = 0x24 << 0
	SRAM4_STALL_DOWNSTREAM      PERFSEL1 = 0x25 << 0
	SRAM4_ACCESS_CONTESTED      PERFSEL1 = 0x26 << 0
	SRAM4_ACCESS                PERFSEL1 = 0x27 << 0
	SRAM3_STALL_UPSTREAM        PERFSEL1 = 0x28 << 0
	SRAM3_STALL_DOWNSTREAM      PERFSEL1 = 0x29 << 0
	SRAM3_ACCESS_CONTESTED      PERFSEL1 = 0x2A << 0
	SRAM3_ACCESS                PERFSEL1 = 0x2B << 0
	SRAM2_STALL_UPSTREAM        PERFSEL1 = 0x2C << 0
	SRAM2_STALL_DOWNSTREAM      PERFSEL1 = 0x2D << 0
	SRAM2_ACCESS_CONTESTED      PERFSEL1 = 0x2E << 0
	SRAM2_ACCESS                PERFSEL1 = 0x2F << 0
	SRAM1_STALL_UPSTREAM        PERFSEL1 = 0x30 << 0
	SRAM1_STALL_DOWNSTREAM      PERFSEL1 = 0x31 << 0
	SRAM1_ACCESS_CONTESTED      PERFSEL1 = 0x32 << 0
	SRAM1_ACCESS                PERFSEL1 = 0x33 << 0
	SRAM0_STALL_UPSTREAM        PERFSEL1 = 0x34 << 0
	SRAM0_STALL_DOWNSTREAM      PERFSEL1 = 0x35 << 0
	SRAM0_ACCESS_CONTESTED      PERFSEL1 = 0x36 << 0
	SRAM0_ACCESS                PERFSEL1 = 0x37 << 0
	XIP_MAIN1_STALL_UPSTREAM    PERFSEL1 = 0x38 << 0
	XIP_MAIN1_STALL_DOWNSTREAM  PERFSEL1 = 0x39 << 0
	XIP_MAIN1_ACCESS_CONTESTED  PERFSEL1 = 0x3A << 0
	XIP_MAIN1_ACCESS            PERFSEL1 = 0x3B << 0
	XIP_MAIN0_STALL_UPSTREAM    PERFSEL1 = 0x3C << 0
	XIP_MAIN0_STALL_DOWNSTREAM  PERFSEL1 = 0x3D << 0
	XIP_MAIN0_ACCESS_CONTESTED  PERFSEL1 = 0x3E << 0
	XIP_MAIN0_ACCESS            PERFSEL1 = 0x3F << 0
	ROM_STALL_UPSTREAM          PERFSEL1 = 0x40 << 0
	ROM_STALL_DOWNSTREAM        PERFSEL1 = 0x41 << 0
	ROM_ACCESS_CONTESTED        PERFSEL1 = 0x42 << 0
	ROM_ACCESS                  PERFSEL1 = 0x43 << 0
)

const (
	PERFSEL1n = 0
)

const (
	PERFSEL2                    PERFSEL2 = 0x7F << 0 //+ Select an event for PERFCTR2. For each downstream port of the main crossbar, four events are available: ACCESS, an access took place; ACCESS_CONTESTED, an access took place that previously stalled due to contention from other masters; STALL_DOWNSTREAM, count cycles where any master stalled due to a stall on the downstream bus; STALL_UPSTREAM, count cycles where any master stalled for any reason, including contention from other masters.
	SIOB_PROC1_STALL_UPSTREAM   PERFSEL2 = 0x00 << 0
	SIOB_PROC1_STALL_DOWNSTREAM PERFSEL2 = 0x01 << 0
	SIOB_PROC1_ACCESS_CONTESTED PERFSEL2 = 0x02 << 0
	SIOB_PROC1_ACCESS           PERFSEL2 = 0x03 << 0
	SIOB_PROC0_STALL_UPSTREAM   PERFSEL2 = 0x04 << 0
	SIOB_PROC0_STALL_DOWNSTREAM PERFSEL2 = 0x05 << 0
	SIOB_PROC0_ACCESS_CONTESTED PERFSEL2 = 0x06 << 0
	SIOB_PROC0_ACCESS           PERFSEL2 = 0x07 << 0
	APB_STALL_UPSTREAM          PERFSEL2 = 0x08 << 0
	APB_STALL_DOWNSTREAM        PERFSEL2 = 0x09 << 0
	APB_ACCESS_CONTESTED        PERFSEL2 = 0x0A << 0
	APB_ACCESS                  PERFSEL2 = 0x0B << 0
	FASTPERI_STALL_UPSTREAM     PERFSEL2 = 0x0C << 0
	FASTPERI_STALL_DOWNSTREAM   PERFSEL2 = 0x0D << 0
	FASTPERI_ACCESS_CONTESTED   PERFSEL2 = 0x0E << 0
	FASTPERI_ACCESS             PERFSEL2 = 0x0F << 0
	SRAM9_STALL_UPSTREAM        PERFSEL2 = 0x10 << 0
	SRAM9_STALL_DOWNSTREAM      PERFSEL2 = 0x11 << 0
	SRAM9_ACCESS_CONTESTED      PERFSEL2 = 0x12 << 0
	SRAM9_ACCESS                PERFSEL2 = 0x13 << 0
	SRAM8_STALL_UPSTREAM        PERFSEL2 = 0x14 << 0
	SRAM8_STALL_DOWNSTREAM      PERFSEL2 = 0x15 << 0
	SRAM8_ACCESS_CONTESTED      PERFSEL2 = 0x16 << 0
	SRAM8_ACCESS                PERFSEL2 = 0x17 << 0
	SRAM7_STALL_UPSTREAM        PERFSEL2 = 0x18 << 0
	SRAM7_STALL_DOWNSTREAM      PERFSEL2 = 0x19 << 0
	SRAM7_ACCESS_CONTESTED      PERFSEL2 = 0x1A << 0
	SRAM7_ACCESS                PERFSEL2 = 0x1B << 0
	SRAM6_STALL_UPSTREAM        PERFSEL2 = 0x1C << 0
	SRAM6_STALL_DOWNSTREAM      PERFSEL2 = 0x1D << 0
	SRAM6_ACCESS_CONTESTED      PERFSEL2 = 0x1E << 0
	SRAM6_ACCESS                PERFSEL2 = 0x1F << 0
	SRAM5_STALL_UPSTREAM        PERFSEL2 = 0x20 << 0
	SRAM5_STALL_DOWNSTREAM      PERFSEL2 = 0x21 << 0
	SRAM5_ACCESS_CONTESTED      PERFSEL2 = 0x22 << 0
	SRAM5_ACCESS                PERFSEL2 = 0x23 << 0
	SRAM4_STALL_UPSTREAM        PERFSEL2 = 0x24 << 0
	SRAM4_STALL_DOWNSTREAM      PERFSEL2 = 0x25 << 0
	SRAM4_ACCESS_CONTESTED      PERFSEL2 = 0x26 << 0
	SRAM4_ACCESS                PERFSEL2 = 0x27 << 0
	SRAM3_STALL_UPSTREAM        PERFSEL2 = 0x28 << 0
	SRAM3_STALL_DOWNSTREAM      PERFSEL2 = 0x29 << 0
	SRAM3_ACCESS_CONTESTED      PERFSEL2 = 0x2A << 0
	SRAM3_ACCESS                PERFSEL2 = 0x2B << 0
	SRAM2_STALL_UPSTREAM        PERFSEL2 = 0x2C << 0
	SRAM2_STALL_DOWNSTREAM      PERFSEL2 = 0x2D << 0
	SRAM2_ACCESS_CONTESTED      PERFSEL2 = 0x2E << 0
	SRAM2_ACCESS                PERFSEL2 = 0x2F << 0
	SRAM1_STALL_UPSTREAM        PERFSEL2 = 0x30 << 0
	SRAM1_STALL_DOWNSTREAM      PERFSEL2 = 0x31 << 0
	SRAM1_ACCESS_CONTESTED      PERFSEL2 = 0x32 << 0
	SRAM1_ACCESS                PERFSEL2 = 0x33 << 0
	SRAM0_STALL_UPSTREAM        PERFSEL2 = 0x34 << 0
	SRAM0_STALL_DOWNSTREAM      PERFSEL2 = 0x35 << 0
	SRAM0_ACCESS_CONTESTED      PERFSEL2 = 0x36 << 0
	SRAM0_ACCESS                PERFSEL2 = 0x37 << 0
	XIP_MAIN1_STALL_UPSTREAM    PERFSEL2 = 0x38 << 0
	XIP_MAIN1_STALL_DOWNSTREAM  PERFSEL2 = 0x39 << 0
	XIP_MAIN1_ACCESS_CONTESTED  PERFSEL2 = 0x3A << 0
	XIP_MAIN1_ACCESS            PERFSEL2 = 0x3B << 0
	XIP_MAIN0_STALL_UPSTREAM    PERFSEL2 = 0x3C << 0
	XIP_MAIN0_STALL_DOWNSTREAM  PERFSEL2 = 0x3D << 0
	XIP_MAIN0_ACCESS_CONTESTED  PERFSEL2 = 0x3E << 0
	XIP_MAIN0_ACCESS            PERFSEL2 = 0x3F << 0
	ROM_STALL_UPSTREAM          PERFSEL2 = 0x40 << 0
	ROM_STALL_DOWNSTREAM        PERFSEL2 = 0x41 << 0
	ROM_ACCESS_CONTESTED        PERFSEL2 = 0x42 << 0
	ROM_ACCESS                  PERFSEL2 = 0x43 << 0
)

const (
	PERFSEL2n = 0
)

const (
	PERFSEL3                    PERFSEL3 = 0x7F << 0 //+ Select an event for PERFCTR3. For each downstream port of the main crossbar, four events are available: ACCESS, an access took place; ACCESS_CONTESTED, an access took place that previously stalled due to contention from other masters; STALL_DOWNSTREAM, count cycles where any master stalled due to a stall on the downstream bus; STALL_UPSTREAM, count cycles where any master stalled for any reason, including contention from other masters.
	SIOB_PROC1_STALL_UPSTREAM   PERFSEL3 = 0x00 << 0
	SIOB_PROC1_STALL_DOWNSTREAM PERFSEL3 = 0x01 << 0
	SIOB_PROC1_ACCESS_CONTESTED PERFSEL3 = 0x02 << 0
	SIOB_PROC1_ACCESS           PERFSEL3 = 0x03 << 0
	SIOB_PROC0_STALL_UPSTREAM   PERFSEL3 = 0x04 << 0
	SIOB_PROC0_STALL_DOWNSTREAM PERFSEL3 = 0x05 << 0
	SIOB_PROC0_ACCESS_CONTESTED PERFSEL3 = 0x06 << 0
	SIOB_PROC0_ACCESS           PERFSEL3 = 0x07 << 0
	APB_STALL_UPSTREAM          PERFSEL3 = 0x08 << 0
	APB_STALL_DOWNSTREAM        PERFSEL3 = 0x09 << 0
	APB_ACCESS_CONTESTED        PERFSEL3 = 0x0A << 0
	APB_ACCESS                  PERFSEL3 = 0x0B << 0
	FASTPERI_STALL_UPSTREAM     PERFSEL3 = 0x0C << 0
	FASTPERI_STALL_DOWNSTREAM   PERFSEL3 = 0x0D << 0
	FASTPERI_ACCESS_CONTESTED   PERFSEL3 = 0x0E << 0
	FASTPERI_ACCESS             PERFSEL3 = 0x0F << 0
	SRAM9_STALL_UPSTREAM        PERFSEL3 = 0x10 << 0
	SRAM9_STALL_DOWNSTREAM      PERFSEL3 = 0x11 << 0
	SRAM9_ACCESS_CONTESTED      PERFSEL3 = 0x12 << 0
	SRAM9_ACCESS                PERFSEL3 = 0x13 << 0
	SRAM8_STALL_UPSTREAM        PERFSEL3 = 0x14 << 0
	SRAM8_STALL_DOWNSTREAM      PERFSEL3 = 0x15 << 0
	SRAM8_ACCESS_CONTESTED      PERFSEL3 = 0x16 << 0
	SRAM8_ACCESS                PERFSEL3 = 0x17 << 0
	SRAM7_STALL_UPSTREAM        PERFSEL3 = 0x18 << 0
	SRAM7_STALL_DOWNSTREAM      PERFSEL3 = 0x19 << 0
	SRAM7_ACCESS_CONTESTED      PERFSEL3 = 0x1A << 0
	SRAM7_ACCESS                PERFSEL3 = 0x1B << 0
	SRAM6_STALL_UPSTREAM        PERFSEL3 = 0x1C << 0
	SRAM6_STALL_DOWNSTREAM      PERFSEL3 = 0x1D << 0
	SRAM6_ACCESS_CONTESTED      PERFSEL3 = 0x1E << 0
	SRAM6_ACCESS                PERFSEL3 = 0x1F << 0
	SRAM5_STALL_UPSTREAM        PERFSEL3 = 0x20 << 0
	SRAM5_STALL_DOWNSTREAM      PERFSEL3 = 0x21 << 0
	SRAM5_ACCESS_CONTESTED      PERFSEL3 = 0x22 << 0
	SRAM5_ACCESS                PERFSEL3 = 0x23 << 0
	SRAM4_STALL_UPSTREAM        PERFSEL3 = 0x24 << 0
	SRAM4_STALL_DOWNSTREAM      PERFSEL3 = 0x25 << 0
	SRAM4_ACCESS_CONTESTED      PERFSEL3 = 0x26 << 0
	SRAM4_ACCESS                PERFSEL3 = 0x27 << 0
	SRAM3_STALL_UPSTREAM        PERFSEL3 = 0x28 << 0
	SRAM3_STALL_DOWNSTREAM      PERFSEL3 = 0x29 << 0
	SRAM3_ACCESS_CONTESTED      PERFSEL3 = 0x2A << 0
	SRAM3_ACCESS                PERFSEL3 = 0x2B << 0
	SRAM2_STALL_UPSTREAM        PERFSEL3 = 0x2C << 0
	SRAM2_STALL_DOWNSTREAM      PERFSEL3 = 0x2D << 0
	SRAM2_ACCESS_CONTESTED      PERFSEL3 = 0x2E << 0
	SRAM2_ACCESS                PERFSEL3 = 0x2F << 0
	SRAM1_STALL_UPSTREAM        PERFSEL3 = 0x30 << 0
	SRAM1_STALL_DOWNSTREAM      PERFSEL3 = 0x31 << 0
	SRAM1_ACCESS_CONTESTED      PERFSEL3 = 0x32 << 0
	SRAM1_ACCESS                PERFSEL3 = 0x33 << 0
	SRAM0_STALL_UPSTREAM        PERFSEL3 = 0x34 << 0
	SRAM0_STALL_DOWNSTREAM      PERFSEL3 = 0x35 << 0
	SRAM0_ACCESS_CONTESTED      PERFSEL3 = 0x36 << 0
	SRAM0_ACCESS                PERFSEL3 = 0x37 << 0
	XIP_MAIN1_STALL_UPSTREAM    PERFSEL3 = 0x38 << 0
	XIP_MAIN1_STALL_DOWNSTREAM  PERFSEL3 = 0x39 << 0
	XIP_MAIN1_ACCESS_CONTESTED  PERFSEL3 = 0x3A << 0
	XIP_MAIN1_ACCESS            PERFSEL3 = 0x3B << 0
	XIP_MAIN0_STALL_UPSTREAM    PERFSEL3 = 0x3C << 0
	XIP_MAIN0_STALL_DOWNSTREAM  PERFSEL3 = 0x3D << 0
	XIP_MAIN0_ACCESS_CONTESTED  PERFSEL3 = 0x3E << 0
	XIP_MAIN0_ACCESS            PERFSEL3 = 0x3F << 0
	ROM_STALL_UPSTREAM          PERFSEL3 = 0x40 << 0
	ROM_STALL_DOWNSTREAM        PERFSEL3 = 0x41 << 0
	ROM_ACCESS_CONTESTED        PERFSEL3 = 0x42 << 0
	ROM_ACCESS                  PERFSEL3 = 0x43 << 0
)

const (
	PERFSEL3n = 0
)
