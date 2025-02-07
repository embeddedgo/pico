// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uart

import "github.com/embeddedgo/pico/hal/system/clock"

type Driver struct {
	p *Periph
}

// NewDriver returns a new driver for p.
func NewDriver(p *Periph) *Driver {
	return &Driver{p: p}
}

// Periph returns the underlying LPSPI peripheral.
func (d *Driver) Periph() *Periph {
	return d.p
}

type Config uint32

const (
	Break      = Config(BRK)        // send break
	ParityEven = Config(PEN | EPS)  // even parity
	ParityOdd  = Config(PEN)        // odd parity
	Stop2b     = Config(STP2)       // two stop bits instead of one
	Word5b     = Config(2 << WLENn) // 5 bit data word
	Word6b     = Config(2 << WLENn) // 6 bit data word
	Word7b     = Config(2 << WLENn) // 7 bit data word
	Word8b     = Config(3 << WLENn) // 8 bit data word
)

func (d *Driver) Config() Config {
	return Config(d.p.LCR_H.Load())
}

func (d *Driver) SetConfig(cfg Config) {
	// TODO: wait for the end of current transfer

	d.p.LCR_H.Store(LCR_H(cfg) | FEN)
}

func (d *Driver) Baudrate() int {
	periHz := clock.PERI.Freq()
	p := d.p
	ibrd := p.IBRD.Load()
	fbrd := p.FBRD.Load()
	return int(4 * periHz / int64(ibrd<<6+fbrd))
}

func (d *Driver) SetBaudrate(baudrate int) {
	// TODO: wait for the end of current transfer

	periHz := clock.PERI.Freq()
	brdiv := uint32(8*periHz/int64(baudrate)) + 1
	ibrd := brdiv >> 7
	fbrd := brdiv & 0x7f >> 1
	if ibrd == 0 {
		ibrd = 1
		fbrd = 0
	} else if ibrd >= 0xffff {
		ibrd = 0xffff
		fbrd = 0
	}
	p := d.p
	p.IBRD.Store(ibrd)
	p.FBRD.Store(fbrd)
	p.LCR_H.Store(p.LCR_H.Load()) // dummy write to latch IBRD and FBRD
}

// Setup resets the driver and the underlying UART peripheral and next calls the
// SetConfig and SetBaudrate methods with the provided arguments.
func (d *Driver) Setup(cfg Config, baudrate int) {
	d.p.SetReset(true)
	d.p.SetReset(false)
	d.SetConfig(cfg)
	d.SetBaudrate(baudrate)
}
