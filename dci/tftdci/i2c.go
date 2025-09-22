// Copyright 2024 The Embedded Go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdci

import (
	"github.com/embeddedgo/device/bus/i2cbus"
	"github.com/embeddedgo/pico/hal/i2c"
)

// I2C is an implementation of the display/tft.DCI interface that uses an I2C
// peripheral to communicate with the I2C displays.
//
// # Limitations
//
// Current implementation is written with the SSD1306 OLED controller in mind
// and the pix/driver/fbdrv driver. It doesn't support ReadBytes, WriteWords,
// WriteWordN. The Cmd and the Write* methods are sending the provided
// commands/data bytes using a single I2C transaction: I2C Start, SSD1306
// control byte, command/data bytes, I2C Stop (the Co bit in the Control Byte is
// cleared).
type I2C struct {
	m    *i2c.Master
	addr i2cbus.Addr
}

// NewI2C returns new I2C based implementation of tftdrv.DCI. User must provide
// a configured I2C master driver and the slave address.
func NewI2C(m *i2c.Master, addr i2cbus.Addr) *I2C {
	return &I2C{m, addr}
}

func (dci *I2C) Driver() *i2c.Master  { return dci.m }
func (dci *I2C) Err(clear bool) error { return dci.m.Err(clear) }

func i2cStart(m *i2c.Master, addr i2cbus.Addr, cmd byte) {
	m.Lock()
	m.SetAddr(addr)
	m.WriteCmd(i2c.Send | int16(cmd))
}

func i2cStop(m *i2c.Master) {
	m.Flush()
	m.Wait(i2c.TX_EMPTY)
	m.Abort() // stop
	m.Unlock()
}

// controll byte values
const (
	singlCmd  = 0b1000_0000
	multiCmd  = 0b0000_0000
	singlData = 0b1100_0000
	multiData = 0b0100_0000
)

func (dci *I2C) Cmd(p []byte, dataMode int) {
	m := dci.m
	i2cStart(m, dci.addr, multiCmd)
	m.WriteBytes(p)
	i2cStop(m)
}

func (dci *I2C) End() {
}

func (dci *I2C) WriteBytes(p []uint8) {
	m := dci.m
	i2cStart(m, dci.addr, multiData)
	m.WriteBytes(p)
	i2cStop(m)
}

func (dci *I2C) WriteString(s string) {
	m := dci.m
	i2cStart(m, dci.addr, multiData)
	m.WriteStr(s)
	i2cStop(m)
}

func (dci *I2C) WriteByteN(b byte, n int) {
	m := dci.m
	i2cStart(m, dci.addr, multiData)
	for n != 0 {
		m.WriteCmd(i2c.Send | int16(b))
		n--
	}
	i2cStop(m)
}

func (dci *I2C) ReadBytes(p []byte) {
	// not supported
}
