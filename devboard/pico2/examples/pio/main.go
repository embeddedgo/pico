// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/device/adc/ads111x"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/pio"
	"github.com/embeddedgo/pico/hal/system/clock"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

const (
	addr = 0b100_1000 // address if the ADDR pin is connected to GND
	cfg  = ads111x.OS | ads111x.AIN0_AIN1 | ads111x.FS2048 | ads111x.SINGLESHOT | ads111x.R8
)

func main() {
	// Used IO pins
	const (
		conTx  = pins.GP0
		conRx  = pins.GP1
		pioOut = pins.GP2
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	pioOut.Setup(iomux.D4mA)
	pioOut.SetAltFunc(iomux.PIO0)

	pb := pio.Block(0)
	pb.SetReset(true)
	pb.SetReset(false)

	pb.INSTR_MEM[0].Store(0xe081) // set pindirs, 1
	pb.INSTR_MEM[1].Store(0xe101) // set pins, 1 [1]
	pb.INSTR_MEM[2].Store(0xe000) // set pins, 0
	pb.INSTR_MEM[3].Store(0x0001) // jmp 1

	sm := &pb.SM[0]
	sm.PINCTRL.StoreBits(pio.SET_BASE, pio.PINCTRL(pioOut)<<pio.SET_BASEn)
	sm.CLKDIV.Store(uint32(clock.PERI.Freq()<<8/1e6) << pio.FRACn)

	pb.CTRL.SetBits(1 << pio.SM_ENABLEn)

	for {
		cfg := pb.DBG_CFGINFO.Load()
		fmt.Println(
			"VERSION:", cfg>>28&0xf,
			"IMEM_SIZE:", cfg>>16&0x3f,
			"SM_COUNT:", cfg>>8&0xf,
			"FIFO_DEPTH:", cfg&0x3f,
		)
		fmt.Println("CLKDIV:", float64(sm.CLKDIV.Load()>>pio.FRACn)/256)
		time.Sleep(2 * time.Second)
	}
}
