// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"sync/atomic"
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
	pb := pio.Block(0)
	pb.SetReset(true)
	pb.SetReset(false)
	pos, _ := pb.Load(pioProg_bt656, 0)

	// Setup the state machine and DMA

	sm := pb.SM(0)
	sm.Configure(pioProg_bt656, pos)
	sm.SetPinBase(advD0, advD0, advD0, advD0)

	const lineWidth = 720

	// Pass two constants to the SM.
	sm.WriteWord32(lineWidth*2 - 1) // line data counter (number of bytes - 1)
	sm.Exec(pio.PULL(false, false, 0))
	sm.Exec(pio.SET(pio.Y, 7, 0)) // SAV XY active field 1 protection bits

	// Double the size of Rx FIFO as we no longer need the Tx one.
	sm.SetFIFOMode(pio.Rx)

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

	const (
		lineWordN  = lineWidth/2 + 1
		logLineNum = 2
		lineNum    = 1 << logLineNum
	)

	var (
		nLineBuf  = make([]uint32, lineWordN*lineNum)
		lineAddrs = allocRingBuf(logLineNum)
	)

	for i := range lineNum {
		p := &nLineBuf[i*lineWordN]
		lineAddrs[i] = unsafe.Pointer(p)
		*p = 0xffff_ffff
	}

	// Alloc two DMA
	dmaData := dma.DMA(0).AllocChannel()
	dmaCtrl := dma.DMA(0).AllocChannel()
	dmaCfg := dma.En | dma.PrioH | dma.S32b

	// The dmaData channel transfers pixels from the PIO to the line buffers.
	dmaData.SetReadAddr(unsafe.Pointer(sm.RxFIFO().Addr()))
	dmaData.SetTransCount(lineWordN, dma.Normal)
	dmaData.SetConfig(dmaCfg|dma.IncW|dma.PIO0_RX0, dmaCtrl)

	// The dmaCtrl triggers the dmaData with the address of the next line buf.
	dmaCtrl.SetReadAddr(unsafe.Pointer(&lineAddrs[0]))
	dmaCtrl.SetWriteAddr(unsafe.Pointer(&dmaData.CTRW()[3]))
	dmaCtrl.SetTransCount(1, dma.Normal)
	dmaCtrl.SetConfigTrig(dmaCfg|dma.IncR|dma.Always|dma.RingSizeCfg(logLineNum+2), dmaCtrl)

	sm.Enable()

	for {
		for i := 0; i < len(nLineBuf); i += lineWordN {
			p := &nLineBuf[i]
		again:
			x := atomic.LoadUint32(p)
			if x == 0xffff_ffff {
				goto again
			}
			*p = 0xffff_ffff
			if x == 0 {
				os.Stdout.WriteString("0")
			} else {
				os.Stdout.WriteString("1")
			}
		}
	}
}

/*
sr := sm.Regs()
fmt.Printf("CLKDIV:    %08x\n", sr.CLKDIV.Load())
fmt.Printf("EXECCTRL:  %08x\n", sr.EXECCTRL.Load())
fmt.Printf("SHIFTCTRL: %08x\n", sr.SHIFTCTRL.Load())
fmt.Printf("ADDR:      %08x\n", sr.ADDR.Load())
fmt.Printf("PINCTRL:   %08x\n", sr.PINCTRL.Load())
*/

func allocRingBuf(logN int) []unsafe.Pointer {
	// There is a chance that the ordinary alignment is OK.
	p := make([]unsafe.Pointer, 1<<logN)
	ptrSize := unsafe.Sizeof(p[0])
	amask := ptrSize<<logN - 1
	if uintptr(unsafe.Pointer(&p[0]))&amask == 0 {
		return p
	}

	// Alloc a bigger one and take an aligned slice from it.
	p = make([]unsafe.Pointer, 2<<logN)
	o := -uintptr(unsafe.Pointer(&p[0])) & amask
	return p[o/ptrSize:][:1<<logN]
}
