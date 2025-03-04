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
	rclk    int
	wclk    int
	started bool
	reconf  bool
}

// NewSPI returns new SPI based implementation of tftdrv.DCI.
func NewSPI(drv *spi.Master, csn, dc iomux.Pin, mode spi.Config, rcHz, wcHz int) *SPI {
	dci := &SPI{
		spi:  drv,
		dc:   gpio.UsePin(dc),
		csn:  gpio.UsePin(csn),
		mode: mode,
		rclk: rcHz,
		wclk: wcHz,
	}

	dci.csn.Set()
	dci.csn.EnableOut()
	dci.dc.Clear()
	dci.dc.EnableOut()
	csn.Setup(iomux.D4mA)
	dc.Setup(iomux.D4mA)

	return dci
}

func (dci *SPI) Driver() *spi.Master  { return dci.spi }
func (dci *SPI) Err(clear bool) error { return nil }

func start(dci *SPI) {
	dci.started = true
	d := dci.spi
	d.Lock()
	d.SetBaudrate(dci.wclk)
	d.Enable()
	dci.csn.Clear()
}

func (dci *SPI) Cmd(p []byte, _ int) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.WaitTxDone()
	dci.dc.Clear()
	d.SetConfig(dci.mode | spi.Word8b)
	d.Write(p)
	d.WaitTxDone()
	dci.dc.Set()
}

// End ends the SPI transaction. It sets CSn pin high, disables the SPI
// peripheral and unlocks the driver. Other usesrs of the same master driver
// can then take controll of the SPI bus locking the driver before use.
func (dci *SPI) End() {
	dci.started = false
	dci.spi.Disable()
	dci.csn.Set()
	dci.spi.Unlock()
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
	d.SetBaudrate(dci.rclk)
	d.Read(p)
	d.SetBaudrate(dci.wclk)
}
