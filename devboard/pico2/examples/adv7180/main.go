// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c0"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/pio"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1

		// ADV data + clock, nine pins: GP2 to GP10
		advD0  = pins.GP2
		advClk = pins.GP10

		advSDA = pins.GP12
		advSCL = pins.GP13
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	// PIO
	for pin := advD0; pin <= advClk; pin++ {
		pin.Setup(iomux.InpEn | iomux.OutDis)
		pin.SetAltFunc(iomux.PIO0)
	}
	pio0 := pio.Block(0)
	pio0.SetReset(true)
	pio0.SetReset(false)
	smDataPos := 0
	smCtrlPos := pioProg_bt656data.Len()
	pio0.Load(pioProg_bt656data, smDataPos)
	pio0.Load(pioProg_bt656ctrl, smCtrlPos)

	// Setup the state machines.
	smData := pio0.SM(0)
	smData.Configure(pioProg_bt656data, smDataPos, smDataPos)
	smData.SetPinBase(advD0, advD0, advD0, advD0)
	smCtrl := pio0.SM(1)
	smCtrl.Configure(pioProg_bt656ctrl, smCtrlPos, smCtrlPos)

	const (
		lineLen    = 720
		lineSize   = lineLen * 2
		lineBufLen = lineLen/2 + 1
		lineBufNum = 3 // the value required by pioProg_bt656ctrl
	)

	lineBufs := make([]uint32, lineBufLen*lineBufNum)
	for i := 0; i < len(lineBufs); i += lineBufLen {
		lineBufs[i] = 0xffff_ffff
	}

	// Pass the required parameters to the state machines.

	smData.WriteWord32(lineSize - 1) // line data counter (no. bytes - 1)
	smData.Exec(pio.PULL(false, false, 0))
	smData.Exec(pio.SET(pio.Y, 7, 0)) // SAV XY active field 1 protection bits
	smData.SetFIFOMode(pio.Rx)        // double the size of Rx FIFO

	smCtrl.WriteWord32(uint32(uintptr(unsafe.Pointer(&lineBufs[0*lineBufLen]))))
	smCtrl.Exec(pio.PULL(false, false, 0))
	smCtrl.Exec(pio.MOV(pio.X, pio.None, pio.OSR, 0))
	smCtrl.WriteWord32(uint32(uintptr(unsafe.Pointer(&lineBufs[1*lineBufLen]))))
	smCtrl.Exec(pio.PULL(false, false, 0))
	smCtrl.Exec(pio.MOV(pio.Y, pio.None, pio.OSR, 0))
	smCtrl.WriteWord32(uint32(uintptr(unsafe.Pointer(&lineBufs[2*lineBufLen]))))
	smCtrl.Exec(pio.PULL(false, false, 0))

	// I2C
	m := i2c0.Master()
	m.UsePin(advSDA, i2c.SDA)
	m.UsePin(advSCL, i2c.SCL)
	m.Setup(100e3)
	advCtrl := m.NewConn(0x21)

	// Disable the automatic free-run mode (blue screen).
	advCtrl.WriteByte(0x0c)
	advCtrl.WriteByte(0)
	err := advCtrl.Close()
	if err != nil {
		fmt.Println("cannot disable ADV free-run:", err)
		time.Sleep(2 * time.Second)
	}

	advPrintStatus(advCtrl)
	time.Sleep(2 * time.Second)
	fmt.Println("Go...")

	// Alloc two DMA channels.
	dmaData := dma.DMA(0).AllocChannel()
	dmaCtrl := dma.DMA(0).AllocChannel()
	dmaCfg := dma.En | dma.PrioH | dma.S32b

	// The dmaData channel transfers pixels from the PIO to the line buffers.
	dmaData.SetReadAddr(unsafe.Pointer(smData.RxFIFO().Addr()))
	dmaData.SetTransCount(lineBufLen, dma.Normal)
	dmaData.SetConfig(dmaCfg|dma.IncW|dma.PIO0_RX0, dmaData)

	// The dmaCtrl triggers the dmaData with the address of the next line buf.
	dmaCtrl.SetReadAddr(unsafe.Pointer(smCtrl.RxFIFO().Addr()))
	dmaCtrl.SetWriteAddr(unsafe.Pointer(&dmaData.CTRW()[3]))
	dmaCtrl.SetTransCount(1, dma.Endless)
	dmaCtrl.SetConfigTrig(dmaCfg|dma.PIO0_RX1, dmaCtrl)

	smCtrl.Enable()
	smData.Enable() // this starts the whole machinery

	last := uint32(0)
	i := 0
	var buf [1]byte
	for {
		p := &lineBufs[i]
	again:
		x := *p
		if x == 0xffff_ffff {
			goto again
		}
		*p = 0xffff_ffff
		if x != last {
			last = x
			buf[0] = '0' + byte(x)
			os.Stdout.Write(buf[:])
		}
		if i += lineBufLen; i >= len(lineBufs) {
			i = 0
		}
	}
}

func printRegs(sm *pio.SM) {
	sr := sm.Regs()
	fmt.Printf("CLKDIV:    %08x\n", sr.CLKDIV.Load())
	fmt.Printf("EXECCTRL:  %08x\n", sr.EXECCTRL.Load())
	fmt.Printf("SHIFTCTRL: %08x\n", sr.SHIFTCTRL.Load())
	fmt.Printf("ADDR:      %08x\n", sr.ADDR.Load())
	fmt.Printf("PINCTRL:   %08x\n", sr.PINCTRL.Load())
}
