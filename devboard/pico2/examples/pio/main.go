// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/pio"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

func main() {
	// Used IO pins
	const (
		conTx   = pins.GP0
		conRx   = pins.GP1
		pioClk  = pins.GP2
		pioData = pins.GP3
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	pioClk.Setup(iomux.D4mA)
	pioClk.SetAltFunc(iomux.PIO0)
	pioData.Setup(iomux.D4mA)
	pioData.SetAltFunc(iomux.PIO0)

	pb := pio.Block(0)
	pb.SetReset(true)
	pb.SetReset(false)

	pos, _ := pb.Load(pioProg_txOnlySPI, 0)

	cfg := pb.DBG_CFGINFO.Load()
	fmt.Println(
		"\nVERSION:", cfg>>28&0xf,
		"IMEM_SIZE:", cfg>>16&0x3f,
		"SM_COUNT:", cfg>>8&0xf,
		"FIFO_DEPTH:", cfg&0x3f,
	)

	sm := &pb.SM[0]

	fmt.Printf("CTRL:      %#x\n", pb.CTRL.Load())
	fmt.Printf("CLKDIV:    %g\n", float64(sm.CLKDIV.Load()>>pio.FRACn)/256)
	fmt.Printf("EXECCTRL:  %#x\n", sm.EXECCTRL.Load())
	fmt.Printf("SHIFTCTRL: %#x\n", sm.SHIFTCTRL.Load())
	fmt.Printf("ADDR:      %#x\n", sm.ADDR.Load())
	fmt.Printf("PINCTRL:   %#x\n", sm.PINCTRL.Load())
	fmt.Println()

	sm.Configure(pioProg_txOnlySPI, pos)
	sm.PINCTRL.StoreBits(
		pio.OUT_BASE|pio.SET_BASE|pio.SIDESET_BASE,
		pio.PINCTRL(pioData)<<pio.OUT_BASEn|
			pio.PINCTRL(pioClk)<<pio.SET_BASEn|
			pio.PINCTRL(pioClk)<<pio.SIDESET_BASEn,
	)

	sm.SetClkFreq(100e3)
	sm.Enable()

	fmt.Printf("CTRL:      %#x\n", pb.CTRL.Load())
	fmt.Printf("CLKDIV:    %g\n", float64(sm.CLKDIV.Load()>>pio.FRACn)/256)
	fmt.Printf("EXECCTRL:  %#x\n", sm.EXECCTRL.Load())
	fmt.Printf("SHIFTCTRL: %#x\n", sm.SHIFTCTRL.Load())
	fmt.Printf("ADDR:      %#x\n", sm.ADDR.Load())
	fmt.Printf("PINCTRL:   %#x\n", sm.PINCTRL.Load())

	txFIFO := &pb.TXF[sm.Num()]
	txFull := pio.FSTAT(1) << (pio.TXFULLn + sm.Num())
	for i := uint32(0); ; i++ {
		//fmt.Printf("%d: FSTAT: %#x\n", i, pb.FSTAT.Load())
		for pb.FSTAT.LoadBits(txFull) != 0 {
		}
		txFIFO.Store(i)
	}
}
