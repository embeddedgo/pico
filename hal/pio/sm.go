// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

import (
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system/clock"
)

type SM struct {
	r SMRegs
}

func (sm *SM) PIO() *PIO {
	addr := uintptr(unsafe.Pointer(sm)) &^ (pioStep - 1)
	return (*PIO)(unsafe.Pointer(addr))
}

func (sm *SM) Num() int {
	addr := uintptr(unsafe.Pointer(sm))
	return int(addr>>2-2) & (numSM - 1)
}

func (sm *SM) Regs() *SMRegs {
	return &sm.r
}

// Disable disables the state machine (stops executing program).
func (sm *SM) Disable() {
	internal.AtomicClear(&sm.PIO().p.CTRL, (1<<SM_ENABLEn)<<uint(sm.Num()))
}

// Enable enables the state machine (starts executing program)
func (sm *SM) Enable() {
	internal.AtomicSet(&sm.PIO().p.CTRL, (1<<SM_ENABLEn)<<uint(sm.Num()))
}

// Reset disables the state machine, resets its internal state and applies the
// default configuration.
func (sm *SM) Reset() {
	pio := sm.PIO()
	smm := CTRL(1) << uint(sm.Num())

	// Disable.
	internal.AtomicClear(&pio.p.CTRL, smm<<SM_ENABLEn)

	// Default config and clearing FIFOs.
	sm.r.CLKDIV.Store(1 << INTn)
	sm.r.EXECCTRL.Store(31 << WRAP_TOPn)
	sm.r.SHIFTCTRL.Store(FJOIN_RX | OUT_SHIFTDIR | IN_SHIFTDIR)
	sm.r.SHIFTCTRL.Store(OUT_SHIFTDIR | IN_SHIFTDIR) // join-disjoin clears FIFOs
	sm.r.PINCTRL.Store(5 << SET_COUNTn)

	// Clear some internal state.
	internal.AtomicSet(&pio.p.CTRL, smm<<SM_RESTARTn)
	internal.AtomicClear(&pio.p.CTRL, smm<<SM_RESTARTn)
	internal.AtomicSet(&pio.p.CTRL, smm<<CLKDIV_RESTARTn)
	internal.AtomicClear(&pio.p.CTRL, smm<<CLKDIV_RESTARTn)
}

// Configure configures the state machine to run the program prog starting from
// the instruction in the memory slot initPC. It doesn't reset the state machine
// before applying the program configuration (see SM.Reset). It doesn't load the
// program to the instruction memory (see PIO.Load).
func (sm *SM) Configure(prog Program, initPC int) {
	if uint(initPC) >= imCap {
		panic("pio: bad initPC")
	}
	prog.AlterSM(sm)
	sm.Exec(JMP(initPC, Always, 0))
}

// SetClkFreq configures the SM to run at the given frequency. It returns the
// actual frequency which may differ from the freq due to rounding. Sea also
// SetClkDiv.
func (sm *SM) SetClkFreq(freq int64) (actual int64) {
	pclk := clock.PERI.Freq()
	div := pclk * 256 / freq
	if div>>24 != 0 {
		return
	}
	sm.r.CLKDIV.Store(uint32(div) << FRACn)
	return pclk * 256 / div
}

// SetClkDiv configures the SM to run at the clock equal to
// clock.PERI.Freq() * 256 / (divInt * 256 + divFrac). See also SetClkFerq.
func (sm *SM) SetClkDiv(divInt, divFrac uint) {
	sm.r.CLKDIV.Store(uint32(divInt<<16+divFrac<<8) << FRACn)
}

// SetPinBase sets the base pin for out, set and sideset operations.
func (sm *SM) SetPinBase(in, out, set, sideset iomux.Pin) {
	gpioBase := int(sm.PIO().p.GPIOBASE.LoadBits(16))
	inBase := PINCTRL(int(in) - gpioBase)
	outBase := PINCTRL(int(out) - gpioBase)
	setBase := PINCTRL(int(set) - gpioBase)
	sidesetBase := PINCTRL(int(sideset) - gpioBase)
	if outBase > 31 || setBase > 31 || sidesetBase > 31 {
		panic("pio: pin out of range")
	}
	sm.r.PINCTRL.StoreBits(
		OUT_BASE|SET_BASE|SIDESET_BASE|IN_BASE,
		outBase<<OUT_BASEn|
			setBase<<SET_BASEn|
			sidesetBase<<SIDESET_BASEn|
			inBase<<IN_BASEn,
	)
}

// Exec provides an instruction to the state machine for immediate execution.
// The instruction is executed even if the state machine is disabled.
func (sm *SM) Exec(instr uint32) {
	sm.r.INSTR.Store(instr)
}

// SetFIFOMode fifoMode
const (
	TxRx   = SHIFTCTRL(0)
	Rx     = FJOIN_RX
	Tx     = FJOIN_TX
	TxPut  = FJOIN_RX_PUT
	TxGet  = FJOIN_RX_GET
	PutGet = FJOIN_RX_PUT | FJOIN_RX_GET
)

// SetFIFOMode sets the FIFO mode to one of: TxRx, Rx, Tx, TxPut, TxGet, PutGet.
func (sm *SM) SetFIFOMode(fifoMode SHIFTCTRL) {
	sm.r.SHIFTCTRL.StoreBits(FJOIN_RX_GET|FJOIN_RX_PUT|FJOIN_TX|FJOIN_RX, fifoMode)
}

func (sm *SM) ReadByte() (b byte, err error) {
	sn := sm.Num()
	rxEmpty := FSTAT(1) << uint(RXEMPTYn+sn)
	pp := &sm.PIO().p
	for pp.FSTAT.LoadBits(rxEmpty) != 0 {
	}
	return byte(pp.RXF[sn].Load()), nil
}

func (sm *SM) Read(p []byte) (n int, err error) {
	sn := sm.Num()
	rxEmpty := FSTAT(1) << uint(RXEMPTYn+sn)
	pp := &sm.PIO().p
	fstat := &pp.FSTAT
	rxf := &pp.RXF[sn]
	for i := range p {
		for fstat.LoadBits(rxEmpty) != 0 {
		}
		p[i] = byte(rxf.Load())
	}
	return len(p), nil
}

func (sm *SM) ReadWord16() (w uint16, err error) {
	sn := sm.Num()
	rxEmpty := FSTAT(1) << uint(RXEMPTYn+sn)
	pp := &sm.PIO().p
	for pp.FSTAT.LoadBits(rxEmpty) != 0 {
	}
	return uint16(pp.RXF[sn].Load()), nil
}

func (sm *SM) Read16(p []uint16) (n int, err error) {
	sn := sm.Num()
	rxEmpty := FSTAT(1) << uint(RXEMPTYn+sn)
	pp := &sm.PIO().p
	fstat := &pp.FSTAT
	rxf := &pp.RXF[sn]
	for i := range p {
		for fstat.LoadBits(rxEmpty) != 0 {
		}
		p[i] = uint16(rxf.Load())
	}
	return len(p), nil
}

func (sm *SM) ReadWord32() (w uint32, err error) {
	sn := sm.Num()
	rxEmpty := FSTAT(1) << uint(RXEMPTYn+sn)
	pp := &sm.PIO().p
	for pp.FSTAT.LoadBits(rxEmpty) != 0 {
	}
	return pp.RXF[sn].Load(), nil
}

func (sm *SM) Read32(p []uint32) (n int, err error) {
	sn := sm.Num()
	rxEmpty := FSTAT(1) << uint(RXEMPTYn+sn)
	pp := &sm.PIO().p
	fstat := &pp.FSTAT
	rxf := &pp.RXF[sn]
	for i := range p {
		for fstat.LoadBits(rxEmpty) != 0 {
		}
		p[i] = rxf.Load()
	}
	return len(p), nil
}

func (sm *SM) WriteWord32(w uint32) error {
	sn := sm.Num()
	txFull := FSTAT(1) << uint(TXFULLn+sn)
	pp := &sm.PIO().p
	for pp.FSTAT.LoadBits(txFull) != 0 {
	}
	pp.TXF[sn].Store(w)
	return nil
}
