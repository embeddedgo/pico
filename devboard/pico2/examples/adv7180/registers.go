// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/embeddedgo/device/bus/i2cbus"
)

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

func advPrintStatus(advCtrl i2cbus.Conn) {
	var i2cBuf [4]byte

	advCtrl.WriteByte(0x10)
	advCtrl.Read(i2cBuf[:4])
	err := advCtrl.Close()
	if err != nil {
		fmt.Println("cannot read ADV status:", err)
		return
	}
	s1 := i2cBuf[0]
	s3 := i2cBuf[3]

	fmt.Print("\nIn lock:                          ", s1>>0&1)
	fmt.Print("\nf_sc locked:                      ", s1>>2&1)
	fmt.Print("\nAGC follows peak white algorithm: ", s1>>3&1)
	ad := "SECAM 525"
	switch s1 >> 4 & 7 {
	case 0:
		ad = "NTSC M/J"
	case 1:
		ad = "NTSC 4.43"
	case 2:
		ad = "PAL M"
	case 3:
		ad = "PAL 60"
	case 4:
		ad = "PAL B/G/H/I/D"
	case 5:
		ad = "SECAM"
	case 6:
		ad = "PAL Combination N"
	}
	fmt.Print("\nResult of autodetection:          ", ad)
	fmt.Print("\nColor kill active:                ", s1>>7&1)
	fmt.Print("\nHorizontal lock indicator:        ", s3>>0&1)
	fmt.Print("\n50 Hz at output:                  ", s3>>2&1)
	fmt.Print("\nBlue screen:                      ", s3>>4&1)
	fmt.Print("\nField length is correct:          ", s3>>5&1)
	fmt.Print("\nInterlaced:                       ", s3>>6&1)
	fmt.Print("\nReliable PAL swinging bursts:     ", s3>>7&1)

	advCtrl.WriteByte(4)
	advCtrl.Read(i2cBuf[:1])
	err = advCtrl.Close()
	if err != nil {
		fmt.Println("cannot read ADV XOC:", err)
		return
	}
	fmt.Printf("\nExtended output control:         %b", i2cBuf[0])

	fmt.Println()
}
