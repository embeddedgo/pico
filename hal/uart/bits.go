// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

// RSR bits
const (
	FE uint32 = 0x01 << 0 //+ Framing error.
	PE uint32 = 0x01 << 1 //+ Parity error.
	BE uint32 = 0x01 << 2 //+ Break error.
	OE uint32 = 0x01 << 3 //+ Overrun error.
)

type FR uint32

// FR bits
const (
	CTSn FR = 0x01 << 0 //+ Clear to send.
	DSRn FR = 0x01 << 1 //+ Data set ready.
	DCDn FR = 0x01 << 2 //+ Data carrier detect.
	BUSY FR = 0x01 << 3 //+ UART busy.
	RXFE FR = 0x01 << 4 //+ Receive FIFO empty.
	TXFF FR = 0x01 << 5 //+ Transmit FIFO full.
	RXFF FR = 0x01 << 6 //+ Receive FIFO full.
	TXFE FR = 0x01 << 7 //+ Transmit FIFO empty.
	RI   FR = 0x01 << 8 //+ Ring indicator.
)

type LCR_H uint32

// LCR_H bits
const (
	BRK  LCR_H = 0x01 << 0 //+ Send break.
	PEN  LCR_H = 0x01 << 1 //+ Parity enable.
	EPS  LCR_H = 0x01 << 2 //+ Even parity select.
	STP2 LCR_H = 0x01 << 3 //+ Two stop bits select.
	FEN  LCR_H = 0x01 << 4 //+ Enable FIFOs.
	WLEN LCR_H = 0x03 << 5 //+ Word length.
	SPS  LCR_H = 0x01 << 7 //+ Stick parity select.
)

// LCR_H bit ofsets
const WLENn = 5

type CR uint32

// CR bits
const (
	UARTEN CR = 0x01 << 0  //+ UART enable.
	SIREN  CR = 0x01 << 1  //+ SIR enable.
	SIRLP  CR = 0x01 << 2  //+ SIR low-power IrDA mode.
	LBE    CR = 0x01 << 7  //+ Loopback enable.
	TXE    CR = 0x01 << 8  //+ Transmit enable.
	RXE    CR = 0x01 << 9  //+ Receive enable.
	DTRn   CR = 0x01 << 10 //+ Data transmit ready.
	RTSn   CR = 0x01 << 11 //+ Request to send.
	OUT1   CR = 0x01 << 12 //+ Complement of the UART Out1 modem status output.
	OUT2   CR = 0x01 << 13 //+ Complement of the UART Out2 modem status output
	RTSEN  CR = 0x01 << 14 //+ RTS hardware flow control enable.
	CTSEN  CR = 0x01 << 15 //+ CTS hardware flow control enable.
)

// IFLS bits
const (
	TXIFLSEL = 0x07 << 0 //+ Transmit interrupt FIFO level select.
	RXIFLSEL = 0x07 << 3 //+ Receive interrupt FIFO level select
)

// IFLS bit offsets
const (
	TXIFLSELn = 0
	RXIFLSELn = 3
)

type INT uint32

const (
	RIMI  INT = 0x01 << 0  //+ nUARTRI modem interrupt.
	CTSMI INT = 0x01 << 1  //+ nUARTCTS modem interrupt.
	DCDMI INT = 0x01 << 2  //+ nUARTDCD modem interrupt.
	DSRMI INT = 0x01 << 3  //+ nUARTDSR modem interrupt.
	RXI   INT = 0x01 << 4  //+ Receive interrupt.
	TXI   INT = 0x01 << 5  //+ Transmit interrupt.
	RTI   INT = 0x01 << 6  //+ Receive timeout interrupt.
	FEI   INT = 0x01 << 7  //+ Framing error interrupt.
	PEI   INT = 0x01 << 8  //+ Parity error interrupt.
	BEI   INT = 0x01 << 9  //+ Break error interrupt.
	OEI   INT = 0x01 << 10 //+ Overrun error interrupt.
)

type DMACR uint32

// DMACR bits
const (
	RXDMAE   DMACR = 0x01 << 0 //+ Receive DMA enable.
	TXDMAE   DMACR = 0x01 << 1 //+ Transmit DMA enable.
	DMAONERR DMACR = 0x01 << 2 //+ DMA on error.
)
