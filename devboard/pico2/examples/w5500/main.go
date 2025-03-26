// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// W5500 tests the communication with the WIZnet W5500 TCP/IP Ethernet
// controller.
package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/spi"
	"github.com/embeddedgo/pico/hal/spi/spi0dma"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		mosi  = pins.GP3
		miso  = pins.GP4
		csn   = pins.GP5
		sck   = pins.GP6
		rst   = pins.GP7
		irq   = pins.GP8
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	// Setup SPI0 driver
	sm := spi0dma.Master()
	sm.UsePin(miso, spi.RXD)
	sm.UsePin(mosi, spi.TXD)
	sm.UsePin(sck, spi.SCK)
	sm.SetConfig(spi.CPOL1 | spi.CPHA1 | spi.Word8b)
	sm.SetBaudrate(1e6)

	d := NewW5500(sm, csn, rst, irq)

	d.Start(CR|WR, 1)
	d.Write([]byte{
		192, 168, 1, 1, // gateway IP address
		255, 255, 255, 0, // subnet mask
		0, 0, 0, 0, 0, 0, // MAC address
		192, 168, 1, 2, // local IP address
	})
	d.End()

	var ip [4]byte

	for {
		d.Start(RD|CR, 0x39)
		chipVer, _ := d.ReadByte()
		d.End()
		d.Start(RD|CR, 0x0f)
		d.Read(ip[:])
		d.End()
		fmt.Printf(
			"Chip ver: %d, IP addr: %d.%d.%d.%d\n",
			chipVer, ip[0], ip[1], ip[2], ip[3],
		)
		time.Sleep(time.Second)
	}
}

// Control byte
const (
	RD = 0 << 2 // read
	WR = 1 << 2 // write

	CR = 0 << 3 // common register

	R  = 1 << 3 // socket register
	TX = 2 << 3 // socket Tx buffer
	RX = 3 << 3 // socket Rx bufer

	S0 = 0 << 5 // socket 0
	S1 = 1 << 5 // socket 1
	S2 = 2 << 5 // socket 2
	S3 = 3 << 5 // socket 3
	S4 = 4 << 5 // socket 4
	S5 = 5 << 5 // socket 5
	S6 = 6 << 5 // socket 6
	S7 = 7 << 5 // socket 7
)

type W5500 struct {
	*spi.Master
	csn gpio.Bit
	rst gpio.Bit
	irq gpio.Bit
	buf [3]byte
}

func NewW5500(sm *spi.Master, csn, rst, irq iomux.Pin) *W5500 {

	d := &W5500{
		Master: sm,
		csn:    gpio.UsePin(csn),
		rst:    gpio.UsePin(rst),
		irq:    gpio.UsePin(irq),
	}
	irq.Setup(iomux.InpEn | iomux.OutDis)
	d.csn.Set()
	d.csn.EnableOut()
	d.rst.Set()
	d.rst.EnableOut()
	csn.Setup(iomux.D4mA)
	rst.Setup(iomux.D4mA)

	return d
}

func (d *W5500) Start(ctrl uint8, addr uint16) {
	d.buf[0] = byte(addr >> 8)
	d.buf[1] = byte(addr)
	d.buf[2] = ctrl
	d.Master.Lock()
	d.Master.Enable()
	d.csn.Clear()
	d.Master.Write(d.buf[:])
}

func (d *W5500) End() {
	d.Master.Disable()
	d.csn.Set()
	d.Master.Unlock()
}
