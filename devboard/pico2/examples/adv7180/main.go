// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/i2c/i2c0"
	"github.com/embeddedgo/pico/hal/system/console/uartcon"
	"github.com/embeddedgo/pico/hal/uart"
	"github.com/embeddedgo/pico/hal/uart/uart0"
)

func main() {
	// Used IO pins
	const (
		conTx = pins.GP0
		conRx = pins.GP1
		sda   = pins.GP20
		scl   = pins.GP21
	)

	// Serial console
	uartcon.Setup(uart0.Driver(), conRx, conTx, uart.Word8b, 115200, "UART0")

	// I2C
	m := i2c0.Master()
	m.UsePin(sda, i2c.SDA)
	m.UsePin(scl, i2c.SCL)
	m.Setup(100e3)

	var buf [4]byte

	c := m.NewConn(0x21)
	/*
		c.WriteByte(0x58)
		c.WriteByte(0x00) // VS/FIELD
		err := c.Close()
		if err != nil {
			fmt.Println(err)
		}
	*/
	for {
		const addr = 0x10
		c.WriteByte(addr)
		c.Read(buf[:])
		err := c.Close()
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		fmt.Println()
		for i := range buf {
			b := buf[i]
			fmt.Printf(
				"%02x %24s: %02x %04b_%04b\n",
				addr+i, regs[addr+i], b, b>>4, b&0xf,
			)
		}
		time.Sleep(10 * time.Second)
	}

}

var regs = [...]string{
	0x00: "Input control",
	0x01: "Video selection",
	0x03: "Output control",
	0x04: "Extended output control",
	0x05: "Reserved",
	0x06: "Reserved",
	0x07: "Autodetect enable",
	0x08: "Contrast",
	0x09: "Reserved",
	0x0A: "Brightness",
	0x0B: "Hue",
	0x0C: "Default Value Y",
	0x0D: "Default Value C",
	0x0E: "ADI Control 1",
	0x0F: "Power management",
	0x10: "Status 1",
	0x11: "IDENT",
	0x12: "Status 2",
	0x13: "Status 3",
	0x14: "Analog clamp control",
	0x15: "Digital Clamp Control 1",
	0x16: "Reserved",
	0x17: "Shaping Filter Control 1",
	0x18: "Shaping Filter Control 2",
	0x19: "Comb filter control",
	0x1D: "ADI Control 2",
	0x27: "Pixel delay control",
	0x2B: "Misc gain control",
	0x2C: "AGC mode control",
	0x2D: "Chroma Gain 1",
	0x2E: "Chroma Gain 2",
	0x2F: "Luma Gain 1",
	0x30: "Luma Gain 2",
	0x31: "VS/FIELD Control 1",
	0x32: "VS/FIELD Control 2",
	0x33: "VS/FIELD Control 3",
	0x34: "HS Position Control 1",
	0x35: "HS Position Control 2",
	0x36: "HS Position Control 3",
	0x37: "Polarity",
	0x38: "NTSC comb control",
	0x39: "PAL comb control",
	0x3A: "ADC control",
	0x3D: "Manual windocontrol",
	0x41: "Resample control",
	0x48: "Gemstar Control 1",
	0x49: "Gemstar Control 2",
	0x4A: "Gemstar Control 3",
	0x4B: "Gemstar Control 4",
	0x4C: "Gemstar control 5",
	0x4D: "CTI DNR Control 1",
	0x4E: "CTI DNR Control 2",
	0x50: "CTI DNR Control 4",
	0x51: "Lock count",
	0x52: "CVBS_TRIM",
	0x58: "VS/FIELD pin control1",
	0x59: "General-purpose outputs2",
	0x8F: "Free-Run Line Length 1",
	0x99: "CCAP 1",
	0x9A: "CCAP 2",
	0x9B: "Letterbox 1",
	0x9C: "Letterbox 2",
	0x9D: "Letterbox 3",
	0xB2: "CRC enable",
	0xC3: "ADC Switch 1",
	0xC4: "ADC Switch 2",
	0xDC: "Letterbox Control 1",
	0xDD: "Letterbox Control 2",
	0xDE: "ST Noise Readback 1",
	0xDF: "ST Noise Readback 2",
	0xE0: "Reserved",
	0xE1: "SD Offset Cb",
	0xE2: "SD Offset Cr",
	0xE3: "SD Saturation Cb",
	0xE4: "SD Saturation Cr",
	0xE5: "NTSC V bit begin",
	0xE6: "NTSC V bit end",
	0xE7: "NTSC F bit toggle",
	0xE8: "PAL V bit begin",
	0xE9: "PAL V bit end",
	0xEA: "PAL F bit toggle",
	0xEB: "Vblank Control 1",
	0xEC: "Vblank Control 2",
	0xF3: "AFE_CONTROL 1",
	0xF4: "Drive strength",
	0xF8: "IF comp control",
	0xF9: "VS mode control",
	0xFB: "Peaking control",
	0xFC: "Coring threshold",
}
