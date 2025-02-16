// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iomux

import "github.com/embeddedgo/pico/hal/internal"

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

type Config uint32

const (
	FastSR   Config = 1 << 0 // Enable fast slew rate
	Schmitt  Config = 1 << 1 // Enable schmitt trigger
	PullDown Config = 1 << 2 // Pull down enable
	PullUp   Config = 1 << 3 // Pull up enable

	Drive Config = 3 << 4 // Drive strength
	D2mA  Config = 0 << 4 // 2 mA
	D4mA  Config = 1 << 4 // 4 mA
	D8mA  Config = 2 << 4 // 8 mA
	D12mA Config = 3 << 4 // 12 mA

	InpEn  Config = 1 << 6 // Input enable
	OutDis Config = 1 << 7 // Output disable
	ISO    Config = 1 << 8 // Pad isolation control.
)

// Config return pin configuration.
//
//go:nosplit
func (p Pin) Config() Config {
	return Config(pb().pad[p].Load())
}

// Setup configures pin.
//
//go:nosplit
func (p Pin) Setup(cfg Config) {
	pb().pad[p].Store(uint32(cfg))
}

// AltFunc represents a mux mode.
type AltFunc uint32

const (
	Func     AltFunc = 0x1F // Selects pin function
	F0       AltFunc = 1
	F1       AltFunc = 2
	F2       AltFunc = 3
	F3       AltFunc = 4
	F4       AltFunc = 5
	F5       AltFunc = 6
	F6       AltFunc = 7
	F7       AltFunc = 8
	F8       AltFunc = 9
	F9       AltFunc = 10
	F10      AltFunc = 11
	F11      AltFunc = 12
	HSTX             = F0
	SPI              = F1
	UART             = F2
	I2C              = F3
	PWM              = F4
	GPIO             = F5
	PIO0             = F6
	PIO1             = F7
	PIO2             = F8
	USB              = F10
	UART_AUX         = F11

	OutOver   AltFunc = 3 << 12 // Peripheral output override
	OutNormal AltFunc = 0 << 12 // normal
	OutInvert AltFunc = 1 << 12 // inverted
	OutLow    AltFunc = 2 << 12 // force low
	OutHigh   AltFunc = 3 << 12 // force high

	OEOver    AltFunc = 3 << 14 // Peripheral output enable override
	OENormal  AltFunc = 0 << 14 // direct
	OEInvert  AltFunc = 1 << 14 // inverted
	OEDisable AltFunc = 2 << 14 // force disabled
	OEEnable  AltFunc = 3 << 14 // force enabled

	InpOver   AltFunc = 3 << 16 // Peripheral input override
	InpNormal AltFunc = 0 << 16 // normal
	InpInvert AltFunc = 1 << 16 // inverted
	InpLow    AltFunc = 2 << 16 // force low
	InpHigh   AltFunc = 3 << 16 // force high

	IRQOver   AltFunc = 3 << 28 // Interrupt override
	IRQNormal AltFunc = 0 << 28 // normal
	IRQInvert AltFunc = 1 << 28 // inverted
	IRQLow    AltFunc = 2 << 28 // force low
	IRQHigh   AltFunc = 3 << 28 // force high
)

// AltFunc returns a currently set muxmode for pin.
//
//go:nosplit
func (p Pin) AltFunc() AltFunc {
	af := AltFunc(ib().gpio[p].ctrl.Load())
	af = af&^Func | (af+1)&Func
	return af
}

// SetAltFunc sets a mux mode for pin.
//
//go:nosplit
func (p Pin) SetAltFunc(af AltFunc) {
	af = af&^Func | (af-1)&Func
	ib().gpio[p].ctrl.Store(uint32(af))
}

// IRQ destination.
const (
	Proc0       int16 = 0 // Processor 0
	Proc1       int16 = 1 // Processor 1
	DormantWake int16 = 2 // Wake ROSC or XOSC from dormant mode
)

// IRQ condition.
const (
	LevelLow  uint8 = 1 << 0 // IRQ when high level
	LevelHigh uint8 = 1 << 1 // IRQ when low level
	EdgeLow   uint8 = 1 << 2 // IRQ when transition from high to low
	EdgeHigh  uint8 = 1 << 3 // IRQ when transition from low to high
)

// SetDstIRQ sets pin as an IRQ source for dst. One pin may be a source for
// multiple destinations with different conditions at the same time.
//
//go:nosplit
func (p Pin) SetDstIRQ(dst int16, condition uint8) {
	i := int(p) >> 3
	shift := 4 * uint(p&7)
	r := &ib().irqCtrl[dst].enable[i]
	internal.AtomicMod(r, 15<<shift, r.Load(), uint32(condition)<<shift)
}

// DstIRQ
//
//go:nosplit
func (p Pin) DstIRQ(dst int16) (condition uint8) {
	i := int(p) >> 3
	shift := 4 * uint(p&7)
	r := &ib().irqCtrl[dst].enable[i]
	return uint8(r.Load() >> shift & 15)
}

// IRQ returns the active IRQ condition for pin.
//
//go:nosplit
func (p Pin) IRQ() (condition uint8) {
	i := int(p) >> 3
	shift := 4 * uint(p&7)
	return uint8(ib().intr[i].Load() >> shift & 15)
}

// IRQ clears the active IRQ condition for pin. Only edge conditions can be
// cleared.
//
//go:nosplit
func (p Pin) ClearIRQ(condition uint8) {
	i := int(p) >> 3
	shift := 4 * uint(p&7)
	ib().intr[i].Store(uint32(condition&15) << shift)
}
