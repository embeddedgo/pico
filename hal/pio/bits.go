// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

type CTRL uint32

const (
	SM_ENABLE               CTRL = 0x0F << 0  //+ Enable/disable each of the four state machines by writing 1/0 to each of these four bits. When disabled, a state machine will cease executing instructions, except those written directly to SMx_INSTR by the system. Multiple bits can be set/cleared at once to run/halt multiple state machines simultaneously.
	SM_RESTART              CTRL = 0x0F << 4  //+ Write 1 to instantly clear internal SM state which may be otherwise difficult to access and will affect future execution. Specifically, the following are cleared: input and output shift counters; the contents of the input shift register; the delay counter; the waiting-on-IRQ state; any stalled instruction written to SMx_INSTR or run by OUT/MOV EXEC; any pin write left asserted due to OUT_STICKY. The contents of the output shift register and the X/Y scratch registers are not affected.
	CLKDIV_RESTART          CTRL = 0x0F << 8  //+ Restart a state machine's clock divider from an initial phase of 0. Clock dividers are free-running, so once started, their output (including fractional jitter) is completely determined by the integer/fractional divisor configured in SMx_CLKDIV. This means that, if multiple clock dividers with the same divisor are restarted simultaneously, by writing multiple 1 bits to this field, the execution clocks of those state machines will run in precise lockstep. Note that setting/clearing SM_ENABLE does not stop the clock divider from running, so once multiple state machines' clocks are synchronised, it is safe to disable/reenable a state machine, whilst keeping the clock dividers in sync. Note also that CLKDIV_RESTART can be written to whilst the state machine is running, and this is useful to resynchronise clock dividers after the divisors (SMx_CLKDIV) have been changed on-the-fly.
	PREV_PIO_MASK           CTRL = 0x0F << 16 //+ A mask of state machines in the neighbouring lower-numbered PIO block in the system (or the highest-numbered PIO block if this is PIO block 0) to which to apply the operations specified by OP_CLKDIV_RESTART, OP_ENABLE, OP_DISABLE in the same write. This allows state machines in a neighbouring PIO block to be started/stopped/clock-synced exactly simultaneously with a write to this PIO block's CTRL register. Neighbouring PIO blocks are disconnected (status signals tied to 0 and control signals ignored) if one block is accessible to NonSecure code, and one is not.
	NEXT_PIO_MASK           CTRL = 0x0F << 20 //+ A mask of state machines in the neighbouring higher-numbered PIO block in the system (or PIO block 0 if this is the highest-numbered PIO block) to which to apply the operations specified by NEXTPREV_CLKDIV_RESTART, NEXTPREV_SM_ENABLE, and NEXTPREV_SM_DISABLE in the same write. This allows state machines in a neighbouring PIO block to be started/stopped/clock-synced exactly simultaneously with a write to this PIO block's CTRL register. Note that in a system with two PIOs, NEXT_PIO_MASK and PREV_PIO_MASK actually indicate the same PIO block. In this case the effects are applied cumulatively (as though the masks were OR'd together). Neighbouring PIO blocks are disconnected (status signals tied to 0 and control signals ignored) if one block is accessible to NonSecure code, and one is not.
	NEXTPREV_SM_ENABLE      CTRL = 0x01 << 24 //+ Write 1 to enable state machines in neighbouring PIO blocks, as specified by NEXT_PIO_MASK and PREV_PIO_MASK in the same write. This is equivalent to setting the corresponding SM_ENABLE bits in those PIOs' CTRL registers. If both OTHERS_SM_ENABLE and OTHERS_SM_DISABLE are set, the disable takes precedence.
	NEXTPREV_SM_DISABLE     CTRL = 0x01 << 25 //+ Write 1 to disable state machines in neighbouring PIO blocks, as specified by NEXT_PIO_MASK and PREV_PIO_MASK in the same write. This is equivalent to clearing the corresponding SM_ENABLE bits in those PIOs' CTRL registers.
	NEXTPREV_CLKDIV_RESTART CTRL = 0x01 << 26 //+ Write 1 to restart the clock dividers of state machines in neighbouring PIO blocks, as specified by NEXT_PIO_MASK and PREV_PIO_MASK in the same write. This is equivalent to writing 1 to the corresponding CLKDIV_RESTART bits in those PIOs' CTRL registers.
)

const (
	SM_ENABLEn               = 0
	SM_RESTARTn              = 4
	CLKDIV_RESTARTn          = 8
	PREV_PIO_MASKn           = 16
	NEXT_PIO_MASKn           = 20
	NEXTPREV_SM_ENABLEn      = 24
	NEXTPREV_SM_DISABLEn     = 25
	NEXTPREV_CLKDIV_RESTARTn = 26
)

type FSTAT uint32

const (
	RXFULL  FSTAT = 0x0F << 0  //+ State machine RX FIFO is full
	RXEMPTY FSTAT = 0x0F << 8  //+ State machine RX FIFO is empty
	TXFULL  FSTAT = 0x0F << 16 //+ State machine TX FIFO is full
	TXEMPTY FSTAT = 0x0F << 24 //+ State machine TX FIFO is empty
)

const (
	RXFULLn  = 0
	RXEMPTYn = 8
	TXFULLn  = 16
	TXEMPTYn = 24
)

type FDEBUG uint32

const (
	RXSTALL FDEBUG = 0x0F << 0  //+ State machine has stalled on full RX FIFO during a blocking PUSH, or an IN with autopush enabled. This flag is also set when a nonblocking PUSH to a full FIFO took place, in which case the state machine has dropped data. Write 1 to clear.
	RXUNDER FDEBUG = 0x0F << 8  //+ RX FIFO underflow (i.e. read-on-empty by the system) has occurred. Write 1 to clear. Note that read-on-empty does not perturb the state of the FIFO in any way, but the data returned by reading from an empty FIFO is undefined, so this flag generally only becomes set due to some kind of software error.
	TXOVER  FDEBUG = 0x0F << 16 //+ TX FIFO overflow (i.e. write-on-full by the system) has occurred. Write 1 to clear. Note that write-on-full does not alter the state or contents of the FIFO in any way, but the data that the system attempted to write is dropped, so if this flag is set, your software has quite likely dropped some data on the floor.
	TXSTALL FDEBUG = 0x0F << 24 //+ State machine has stalled on empty TX FIFO during a blocking PULL, or an OUT with autopull enabled. Write 1 to clear.
)

const (
	RXSTALLn = 0
	RXUNDERn = 8
	TXOVERn  = 16
	TXSTALLn = 24
)

type FLEVEL uint32

const (
	TX0 FLEVEL = 0x0F << 0  //+
	RX0 FLEVEL = 0x0F << 4  //+
	TX1 FLEVEL = 0x0F << 8  //+
	RX1 FLEVEL = 0x0F << 12 //+
	TX2 FLEVEL = 0x0F << 16 //+
	RX2 FLEVEL = 0x0F << 20 //+
	TX3 FLEVEL = 0x0F << 24 //+
	RX3 FLEVEL = 0x0F << 28 //+
)

const (
	TX0n = 0
	RX0n = 4
	TX1n = 8
	RX1n = 12
	TX2n = 16
	RX2n = 20
	TX3n = 24
	RX3n = 28
)

type DBG_CFGINFO uint32

const (
	FIFO_DEPTH DBG_CFGINFO = 0x3F << 0  //+ The depth of the state machine TX/RX FIFOs, measured in words. Joining fifos via SHIFTCTRL_FJOIN gives one FIFO with double this depth.
	SM_COUNT   DBG_CFGINFO = 0x0F << 8  //+ The number of state machines this PIO instance is equipped with.
	IMEM_SIZE  DBG_CFGINFO = 0x3F << 16 //+ The size of the instruction memory, measured in units of one instruction
	VERSION    DBG_CFGINFO = 0x0F << 28 //+ Version of the core PIO hardware.
	V0         DBG_CFGINFO = 0x00 << 28 //  Version 0 (RP2040)
	V1         DBG_CFGINFO = 0x01 << 28 //  Version 1 (RP2350)
)

const (
	FIFO_DEPTHn = 0
	SM_COUNTn   = 8
	IMEM_SIZEn  = 16
	VERSIONn    = 28
)

// CLKDIV
const (
	FRAC uint32 = 0xFF << 8    //+ Fractional part of clock divisor
	INT  uint32 = 0xFFFF << 16 //+ Effective frequency is sysclk/(int + frac/256). Value of 0 is interpreted as 65536. If INT is 0, FRAC must also be 0.
)

const (
	FRACn = 8
	INTn  = 16
)

type EXECCTRL uint32

const (
	STATUS_N      EXECCTRL = 0x1F << 0  //+ Comparison level or IRQ index for the MOV x, STATUS instruction. If STATUS_SEL is TXLEVEL or RXLEVEL, then values of STATUS_N greater than the current FIFO depth are reserved, and have undefined behaviour.
	IRQ_THISPIO   EXECCTRL = 0x00 << 0  //  Index 0-7 of an IRQ flag in this PIO block
	IRQ_PREVPIO   EXECCTRL = 0x08 << 0  //  Index 0-7 of an IRQ flag in the next lower-numbered PIO block
	IRQ_NEXTPIO   EXECCTRL = 0x10 << 0  //  Index 0-7 of an IRQ flag in the next higher-numbered PIO block
	STATUS_SEL    EXECCTRL = 0x03 << 5  //+ Comparison used for the MOV x, STATUS instruction.
	TXLEVEL       EXECCTRL = 0x00 << 5  //  All-ones if TX FIFO level < N, otherwise all-zeroes
	RXLEVEL       EXECCTRL = 0x01 << 5  //  All-ones if RX FIFO level < N, otherwise all-zeroes
	IRQ           EXECCTRL = 0x02 << 5  //  All-ones if the indexed IRQ flag is raised, otherwise all-zeroes
	WRAP_BOTTOM   EXECCTRL = 0x1F << 7  //+ After reaching wrap_top, execution is wrapped to this address.
	WRAP_TOP      EXECCTRL = 0x1F << 12 //+ After reaching this address, execution is wrapped to wrap_bottom. If the instruction is a jump, and the jump condition is true, the jump takes priority.
	OUT_STICKY    EXECCTRL = 0x01 << 17 //+ Continuously assert the most recent OUT/SET to the pins
	INLINE_OUT_EN EXECCTRL = 0x01 << 18 //+ If 1, use a bit of OUT data as an auxiliary write enable When used in conjunction with OUT_STICKY, writes with an enable of 0 will deassert the latest pin write. This can create useful masking/override behaviour due to the priority ordering of state machine pin writes (SM0 < SM1 < ...)
	OUT_EN_SEL    EXECCTRL = 0x1F << 19 //+ Which data bit to use for inline OUT enable
	JMP_PIN       EXECCTRL = 0x1F << 24 //+ The GPIO number to use as condition for JMP PIN. Unaffected by input mapping.
	SIDE_PINDIR   EXECCTRL = 0x01 << 29 //+ If 1, side-set data is asserted to pin directions, instead of pin values
	SIDE_EN       EXECCTRL = 0x01 << 30 //+ If 1, the MSB of the Delay/Side-set instruction field is used as side-set enable, rather than a side-set data bit. This allows instructions to perform side-set optionally, rather than on every instruction, but the maximum possible side-set width is reduced from 5 to 4. Note that the value of PINCTRL_SIDESET_COUNT is inclusive of this enable bit.
	EXEC_STALLED  EXECCTRL = 0x01 << 31 //+ If 1, an instruction written to SMx_INSTR is stalled, and latched by the state machine. Will clear to 0 once this instruction completes.
)

const (
	STATUS_Nn      = 0
	STATUS_SELn    = 5
	WRAP_BOTTOMn   = 7
	WRAP_TOPn      = 12
	OUT_STICKYn    = 17
	INLINE_OUT_ENn = 18
	OUT_EN_SELn    = 19
	JMP_PINn       = 24
	SIDE_PINDIRn   = 29
	SIDE_ENn       = 30
	EXEC_STALLEDn  = 31
)

type SHIFTCTRL uint32

const (
	IN_COUNT     SHIFTCTRL = 0x1F << 0  //+ Set the number of pins which are not masked to 0 when read by an IN PINS, WAIT PIN or MOV x, PINS instruction. For example, an IN_COUNT of 5 means that the 5 LSBs of the IN pin group are visible (bits 4:0), but the remaining 27 MSBs are masked to 0. A count of 32 is encoded with a field value of 0, so the default behaviour is to not perform any masking. Note this masking is applied in addition to the masking usually performed by the IN instruction. This is mainly useful for the MOV x, PINS instruction, which otherwise has no way of masking pins.
	FJOIN_RX_GET SHIFTCTRL = 0x01 << 14 //+ If 1, disable this state machine's RX FIFO, make its storage available for random read access by the state machine (using the `get` instruction) and, unless FJOIN_RX_PUT is also set, random write access by the processor (through the RXFx_PUTGETy registers). If FJOIN_RX_PUT and FJOIN_RX_GET are both set, then the RX FIFO's registers can be randomly read/written by the state machine, but are completely inaccessible to the processor. Setting this bit will clear the FJOIN_TX and FJOIN_RX bits.
	FJOIN_RX_PUT SHIFTCTRL = 0x01 << 15 //+ If 1, disable this state machine's RX FIFO, make its storage available for random write access by the state machine (using the `put` instruction) and, unless FJOIN_RX_GET is also set, random read access by the processor (through the RXFx_PUTGETy registers). If FJOIN_RX_PUT and FJOIN_RX_GET are both set, then the RX FIFO's registers can be randomly read/written by the state machine, but are completely inaccessible to the processor. Setting this bit will clear the FJOIN_TX and FJOIN_RX bits.
	AUTOPUSH     SHIFTCTRL = 0x01 << 16 //+ Push automatically when the input shift register is filled, i.e. on an IN instruction which causes the input shift counter to reach or exceed PUSH_THRESH.
	AUTOPULL     SHIFTCTRL = 0x01 << 17 //+ Pull automatically when the output shift register is emptied, i.e. on or following an OUT instruction which causes the output shift counter to reach or exceed PULL_THRESH.
	IN_SHIFTDIR  SHIFTCTRL = 0x01 << 18 //+ 1 = shift input shift register to right (data enters from left). 0 = to left.
	OUT_SHIFTDIR SHIFTCTRL = 0x01 << 19 //+ 1 = shift out of output shift register to right. 0 = to left.
	PUSH_THRESH  SHIFTCTRL = 0x1F << 20 //+ Number of bits shifted into ISR before autopush, or conditional push (PUSH IFFULL), will take place. Write 0 for value of 32.
	PULL_THRESH  SHIFTCTRL = 0x1F << 25 //+ Number of bits shifted out of OSR before autopull, or conditional pull (PULL IFEMPTY), will take place. Write 0 for value of 32.
	FJOIN_TX     SHIFTCTRL = 0x01 << 30 //+ When 1, TX FIFO steals the RX FIFO's storage, and becomes twice as deep. RX FIFO is disabled as a result (always reads as both full and empty). FIFOs are flushed when this bit is changed.
	FJOIN_RX     SHIFTCTRL = 0x01 << 31 //+ When 1, RX FIFO steals the TX FIFO's storage, and becomes twice as deep. TX FIFO is disabled as a result (always reads as both full and empty). FIFOs are flushed when this bit is changed.
)

const (
	IN_COUNTn     = 0
	FJOIN_RX_GETn = 14
	FJOIN_RX_PUTn = 15
	AUTOPUSHn     = 16
	AUTOPULLn     = 17
	IN_SHIFTDIRn  = 18
	OUT_SHIFTDIRn = 19
	PUSH_THRESHn  = 20
	PULL_THRESHn  = 25
	FJOIN_TXn     = 30
	FJOIN_RXn     = 31
)

type PINCTRL uint32

const (
	OUT_BASE      PINCTRL = 0x1F << 0  //+ The lowest-numbered pin that will be affected by an OUT PINS, OUT PINDIRS or MOV PINS instruction. The data written to this pin will always be the least-significant bit of the OUT or MOV data.
	SET_BASE      PINCTRL = 0x1F << 5  //+ The lowest-numbered pin that will be affected by a SET PINS or SET PINDIRS instruction. The data written to this pin is the least-significant bit of the SET data.
	SIDESET_BASE  PINCTRL = 0x1F << 10 //+ The lowest-numbered pin that will be affected by a side-set operation. The MSBs of an instruction's side-set/delay field (up to 5, determined by SIDESET_COUNT) are used for side-set data, with the remaining LSBs used for delay. The least-significant bit of the side-set portion is the bit written to this pin, with more-significant bits written to higher-numbered pins.
	IN_BASE       PINCTRL = 0x1F << 15 //+ The pin which is mapped to the least-significant bit of a state machine's IN data bus. Higher-numbered pins are mapped to consecutively more-significant data bits, with a modulo of 32 applied to pin number.
	OUT_COUNT     PINCTRL = 0x3F << 20 //+ The number of pins asserted by an OUT PINS, OUT PINDIRS or MOV PINS instruction. In the range 0 to 32 inclusive.
	SET_COUNT     PINCTRL = 0x07 << 26 //+ The number of pins asserted by a SET. In the range 0 to 5 inclusive.
	SIDESET_COUNT PINCTRL = 0x07 << 29 //+ The number of MSBs of the Delay/Side-set instruction field which are used for side-set. Inclusive of the enable bit, if present. Minimum of 0 (all delay bits, no side-set) and maximum of 5 (all side-set, no delay).
)

const (
	OUT_BASEn      = 0
	SET_BASEn      = 5
	SIDESET_BASEn  = 10
	IN_BASEn       = 15
	OUT_COUNTn     = 20
	SET_COUNTn     = 26
	SIDESET_COUNTn = 29
)

type INTR uint32

const (
	SM0_RXNEMPTY INTR = 0x01 << 0  //+
	SM1_RXNEMPTY INTR = 0x01 << 1  //+
	SM2_RXNEMPTY INTR = 0x01 << 2  //+
	SM3_RXNEMPTY INTR = 0x01 << 3  //+
	SM0_TXNFULL  INTR = 0x01 << 4  //+
	SM1_TXNFULL  INTR = 0x01 << 5  //+
	SM2_TXNFULL  INTR = 0x01 << 6  //+
	SM3_TXNFULL  INTR = 0x01 << 7  //+
	SM0          INTR = 0x01 << 8  //+
	SM1          INTR = 0x01 << 9  //+
	SM2          INTR = 0x01 << 10 //+
	SM3          INTR = 0x01 << 11 //+
	SM4          INTR = 0x01 << 12 //+
	SM5          INTR = 0x01 << 13 //+
	SM6          INTR = 0x01 << 14 //+
	SM7          INTR = 0x01 << 15 //+
)

const (
	SM0_RXNEMPTYn = 0
	SM1_RXNEMPTYn = 1
	SM2_RXNEMPTYn = 2
	SM3_RXNEMPTYn = 3
	SM0_TXNFULLn  = 4
	SM1_TXNFULLn  = 5
	SM2_TXNFULLn  = 6
	SM3_TXNFULLn  = 7
	SM0n          = 8
	SM1n          = 9
	SM2n          = 10
	SM3n          = 11
	SM4n          = 12
	SM5n          = 13
	SM6n          = 14
	SM7n          = 15
)
