// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/embeddedgo/display/pix/driver/tftdrv"
	"github.com/embeddedgo/display/pix/driver/tftdrv/st7789"
	"github.com/embeddedgo/pico/dci/tftdci"
	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/devboard/pico2/module/waveshare/pico-restouch-lcd-2.8/lcd"
	"github.com/embeddedgo/pico/hal/dma"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c1"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/pio"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart1"
)

func main() {
	// Used IO pins
	const (
		// ADV data + clock, nine pins: GP0 to GP7 and GP14
		advD0  = pins.GP0
		advD7  = pins.GP7
		advClk = pins.GP14

		// Serial console
		conTx = pins.GP20
		conRx = pins.GP21

		// ADV I2C control
		advSDA = pins.GP26_A0
		advSCL = pins.GP27_A1
	)

	// Serial console
	uartcon.Setup(uart1.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	// PIO
	for pin := advD0; pin <= advD7; pin++ {
		pin.Setup(iomux.InpEn | iomux.OutDis)
		pin.SetAltFunc(iomux.PIO0)
	}
	advClk.Setup(iomux.InpEn | iomux.OutDis)
	advClk.SetAltFunc(iomux.PIO0)

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
	m := i2c1.Master()
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

	var cxy [4]byte
	disp := lcd.Display
	disp.SetDir(-1)

	// To speed things up use the tftdrv.DCI directly.
	dci := disp.Driver().(*tftdrv.DriverOver).DCI().(*tftdci.SPI)

	// Set the drawing window to the whole display
	siz := disp.Bounds().Size()
	width := min(siz.X, 320)
	height := siz.Y
	cxy[0] = st7789.CASET
	dci.Cmd(cxy[:1], tftdrv.Write)
	cxy[0] = 0
	cxy[2] = uint8(width >> 8)
	cxy[3] = uint8(width)
	dci.WriteBytes(cxy[:])
	cxy[0] = st7789.PASET
	dci.Cmd(cxy[:1], tftdrv.Write)
	cxy[0] = 0
	cxy[2] = uint8(height >> 8)
	cxy[3] = uint8(height)
	dci.WriteBytes(cxy[:])
	cxy[0] = st7789.RAMWR
	dci.Cmd(cxy[:1], tftdrv.Write)
	dci.WriteWordN(0xffff, width*height)

	time.Sleep(2 * time.Second)
	fmt.Println("Go...")

	ch := make(chan []uint16, 2)
	buf := make([]uint16, 3*lineBufLen)
	go writeLine(dci, ch)

	i := 0
	for {
		p := &lineBufs[i]
	again:
		x := *p
		if x == 0xffff_ffff {
			goto again
		}
		*p = 0xffff_ffff
		src := lineBufs[i+1 : i+(width+1)]
		dst := buf[i+1 : i+(width+1)]
		for k, w := range src {
			y := w >> 10 & 63
			dst[k] = uint16(y << 5)
		}
		ch <- dst

		if i += lineBufLen; i >= len(lineBufs) {
			i = 0
		}
	}
}

func writeLine(dci *tftdci.SPI, ch <-chan []uint16) {
	for p := range ch {
		dci.WriteWords(p)
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
