// Copyright 2024 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package i2c

import (
	"unsafe"

	"github.com/embeddedgo/device/bus/i2cbus"
)

// Name implements the i2cbus.Master interface. The default name is the name of
// the underlying peripheral (e.g. "I2C0") but can be changed using SetName.
func (d *Master) Name() string {
	return d.name
}

// SetName allows to change the default master name (see Name).
func (d *Master) SetName(s string) {
	d.name = s
}

// SetID sets the Master ID. Its three least significant bits are used for
// arbitration between competing masters while switching to the High Speed mode
// (not supported by RP2350).
func (d *Master) SetID(id uint8) {
	d.id = id
}

// ID rteturns the Master ID. See SetID for more information.
func (d *Master) ID() uint8 {
	return d.id
}

type conn struct {
	d    *Master
	a    i2cbus.Addr
	open bool
}

// NewConn implements the i2cbus.Master interface.
func (d *Master) NewConn(a i2cbus.Addr) i2cbus.Conn {
	return &conn{d: d, a: a}
}

// Addr implements the i2cbus.Conn interface.
func (c *conn) Addr() i2cbus.Addr {
	return c.a
}

// Master implements the i2cbus.Conn interface.
func (c *conn) Master() i2cbus.Master {
	return c.d
}

func open(c *conn) {
	c.open = true
	d := c.d
	d.Lock()
	if d.p.TAR.Load() != TAR(c.a)&0x3ff {
		d.SetAddr(c.a)
	}
}

// Write implements the i2cbus.Conn interface and the io.Writer interface.
func (c *conn) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	if !c.open {
		open(c)
	}
	d := c.d
	d.WriteBytes(p)
	d.Flush() // ensure p isn't used after return
	err = d.Err(false)
	if err == nil {
		n = len(p)
	}
	return
}

// WriteString implements the io.StringWriter interface.
func (c *conn) WriteString(s string) (n int, err error) {
	return c.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// WriteByte implements the i2cbus.Conn interface and the io.ByteWriter
// interface.
func (c *conn) WriteByte(b byte) error {
	if !c.open {
		open(c)
	}
	d := c.d
	d.WriteCmd(Send | int16(b))
	return d.Err(false)
}

// Read implements the i2cbus.Conn interface and the io.Reader interface.
func (c *conn) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		return
	}
	if n > 256 {
		n = 256
	}
	if !c.open {
		open(c)
	}
	d := c.d
	d.WriteCmd(Recv | int16(n-1))
	d.ReadBytes(p)
	err = d.Err(false)
	if err != nil {
		n = 0
	}
	return
}

// ReadByte implements the i2cbus.Conn interface and the io.ByteReader
// interface.
func (c *conn) ReadByte() (b byte, err error) {
	if !c.open {
		open(c)
	}
	d := c.d
	d.WriteCmd(Recv)
	b = d.ReadByte()
	err = d.Err(false)
	return
}

// Close implements the i2cbus.Conn interface and the io.Closer interface.
func (c *conn) Close() error {
	if !c.open {
		return nil // already closed
	}
	d := c.d
	d.Wait(TX_EMPTY)
	d.Abort() // stop
	err := d.Err(true)
	d.Unlock()
	c.open = false
	return err
}
