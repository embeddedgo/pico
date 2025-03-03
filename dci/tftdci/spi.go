// Copyright 2025 The Embedded Go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdci

import (
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/spi"
)

// An SPI is an implementation of the display/tft.DCI interface that uses an SPI
// peripheral to communicate with the display in what is known as 4-line mode.
type SPI struct {
	spi     *spi.Master
	dc      gpio.Bit
	csn     gpio.Bit
	mode    spi.Config
	rclkHz  int
	wclkHz  int
	started bool
	reconf  bool
}

// NewSPI returns new SPI based implementation of tftdrv.DCI.
func NewSPI(drv *spi.Master, csn, dc iomux.Pin, mode spi.Config, rclkHz, wclkHz int) *SPI {
	dci := new(SPI)
	dci.spi = drv
	dci.dc = gpio.BitForPin(dc)
	dci.csn = gpio.BitForPin(csn)
	dci.mode = mode

	gpio.UsePin(csn)
	gpio.UsePin(dc)
	dci.csn.Set()
	dci.dc.Clear()
	csn.Setup(iomux.D4mA)
	dc.Setup(iomux.D4mA)

	return dci
}

func (dci *SPI) Driver() *spi.Master  { return dci.spi }
func (dci *SPI) Err(clear bool) error { return nil }

func start(dci *SPI) {
	dci.started = true
	d := dci.spi
	d.SetBaudrate(dci.wclkHz)
	d.Enable()
	dci.csn.Clear()
}

func (dci *SPI) Cmd(p []byte, _ int) {
	if !dci.started {
		start(dci)
	}
	dci.dc.Clear()
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word8b)
	d.Write(p)
	dci.dc.Set()
}

// End ends the SPI transaction. It sets CSn pin high and disables the SPI
// peripheral. Other usesrs of the same master driver can then take controll of
// the SPI bus.
func (dci *SPI) End() {
	dci.started = false
	dci.csn.Set()
	dci.spi.Disable()
}

func (dci *SPI) WriteBytes(p []uint8) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word8b)
	d.Write(p)
}

func (dci *SPI) WriteString(s string) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word8b)
	d.WriteString(s)
}

func (dci *SPI) WriteByteN(b byte, n int) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word8b)
	d.WriteByteN(b, n)
}

func (dci *SPI) WriteWords(p []uint16) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word16b)
	d.Write16(p)
}

func (dci *SPI) WriteWordN(w uint16, n int) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word16b)
	d.WriteWord16N(w, n)
}

func (dci *SPI) ReadBytes(p []byte) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.SetConfig(dci.mode | spi.Word8b)
	d.SetBaudrate(dci.rclkHz)
	d.Read(p)
	d.SetBaudrate(dci.wclkHz)
}
