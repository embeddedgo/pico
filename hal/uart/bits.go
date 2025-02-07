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
	CTS  FR = 0x01 << 0 //+ Clear to send.
	DSR  FR = 0x01 << 1 //+ Data set ready.
	DCD  FR = 0x01 << 2 //+ Data carrier detect.
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
	DTR    CR = 0x01 << 10 //+ Data transmit ready.
	RTS    CR = 0x01 << 11 //+ Request to send.
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

type IS uint32

const (
	RIRMIS  IS = 0x01 << 0  //+ nUARTRI modem interrupt status.
	CTSRMIS IS = 0x01 << 1  //+ nUARTCTS modem interrupt status.
	DCDRMIS IS = 0x01 << 2  //+ nUARTDCD modem interrupt status.
	DSRRMIS IS = 0x01 << 3  //+ nUARTDSR modem interrupt status.
	RXRIS   IS = 0x01 << 4  //+ Receive interrupt status.
	TXRIS   IS = 0x01 << 5  //+ Transmit interrupt status.
	RTRIS   IS = 0x01 << 6  //+ Receive timeout interrupt status.
	FERIS   IS = 0x01 << 7  //+ Framing error interrupt status.
	PERIS   IS = 0x01 << 8  //+ Parity error interrupt status.
	BERIS   IS = 0x01 << 9  //+ Break error interrupt status.
	OERIS   IS = 0x01 << 10 //+ Overrun error interrupt status.
)

type DMACR uint32

// DMACR bits
const (
	RXDMAE   DMACR = 0x01 << 0 //+ Receive DMA enable.
	TXDMAE   DMACR = 0x01 << 1 //+ Transmit DMA enable.
	DMAONERR DMACR = 0x01 << 2 //+ DMA on error.
)
