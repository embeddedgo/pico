// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iomux

// Pin represents an I/O pin (pad).
type Pin int16

const (
	P00 Pin = iota
	P01
	P02
	P03
	P04
	P05
	P06
	P07
	P08
	P09
	P10
	P11
	P12
	P13
	P14
	P15
	P16
	P17
	P18
	P19
	P20
	P21
	P22
	P23
	P24
	P25
	P26
	P27
	P28
	P29
	P30
	P31
	P32
	P33
	P34
	P35
	P36
	P37
	P38
	P39
	P40
	P41
	P42
	P43
	P44
	P45
	P46
	P47
	SWCLK
	SWD
)

/*
type Config uint32

const (
	FastSR Config = 1 << 0 // Enable fast slew rate

	Drive  Config = 7 << 3 // Drive strength field
	Drive0 Config = 0 << 3 // Rout = ∞Ω (output driver disabled)
	Drive1 Config = 1 << 3 // Rout = R, R = 150Ω @ 3V3, R = 260Ω @ 1V8
	Drive2 Config = 2 << 3 // Rout = R / 2
	Drive3 Config = 3 << 3 // Rout = R / 3
	Drive4 Config = 4 << 3 // Rout = R / 4
	Drive5 Config = 5 << 3 // Rout = R / 5
	Drive6 Config = 6 << 3 // Rout = R / 6
	Drive7 Config = 7 << 3 // Rout = R / 7

	Speed       Config = 3 << 6 // Speed field
	SpeedLow    Config = 0 << 6 // Speed low (50MHz)
	SpeedMedium Config = 1 << 6 // Speed medium (100MHz)
	SpeedFast   Config = 2 << 6 // Speed fast (150MHz)
	SpeedMax    Config = 3 << 6 // Speed max (200MHz)

	OpenDrain Config = 1 << 11 // Enable open drain mode

	PK   Config = 3 << 12 // Pull / keep field
	Keep Config = 1 << 12 // Enable pull/keep mode
	Pull Config = 3 << 12 // Use pull mode instead of keep mode

	PullSel  Config = 3 << 14 // Select pull direction and strength
	Down100k Config = 0 << 14 // 100kΩ pull-down
	Up47k    Config = 1 << 14 //  47kΩ pull-up
	Up100k   Config = 2 << 14 // 100kΩ pull up
	Up22k    Config = 3 << 14 //  22kΩ pull up

	Hys Config = 1 << 16 //+ Enable hysteresis mode
)

// AltFunc represents a mux mode.
type AltFunc int8

const (
	ALT   AltFunc = 0xf << 0 // Mux mode select field
	ALT0  AltFunc = 0x0 << 0 // Select ALT0 mux mode
	ALT1  AltFunc = 0x1 << 0 // Select ALT1 mux mode
	ALT2  AltFunc = 0x2 << 0 // Select ALT2 mux mode
	ALT3  AltFunc = 0x3 << 0 // Select ALT3 mux mode
	ALT4  AltFunc = 0x4 << 0 // Select ALT4 mux mode
	ALT5  AltFunc = 0x5 << 0 // Select ALT5 mux mode
	ALT6  AltFunc = 0x6 << 0 // Select ALT6 mux mode
	ALT7  AltFunc = 0x7 << 0 // Select ALT7 mux mode
	ALT8  AltFunc = 0x8 << 0 // Select ALT8 mux mode
	ALT9  AltFunc = 0x9 << 0 // Select ALT9 mux mode
	ALT10 AltFunc = 0xa << 0 // Select ALT10 mux mode

	SION AltFunc = 0x1 << 4 // Software Input On field

	GPIO = ALT5 // More readable alias for ALT5
)

// Config return pin configuration.
func (p Pin) Config() Config {
	return Config(pr().pad[p].Load())
}

// Setup configures pin.
func (p Pin) Setup(cfg Config) {
	pr().pad[p].Store(uint32(cfg))
}

// AltFunc returns a currently set muxmode for pin.
func (p Pin) AltFunc() AltFunc {
	return AltFunc(pr().mux[p].Load())
}

// SetAltFunc sets a mux mode for pin.
func (p Pin) SetAltFunc(af AltFunc) {
	pr().mux[p].Store(uint32(af))
}
*/