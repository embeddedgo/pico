// Copyright 2025 The Embedded Go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftdci

import (
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/spi"
)

// An SPI is an implementation of the tftdrv.DCI interface that uses an SPI
// peripheral to communicate with the display using the 4-line Serial Protocol.
type SPI struct {
	spi     *spi.Master
	dc      gpio.Bit
	csn     gpio.Bit
	mode8   spi.Config
	mode16  spi.Config
	rclk    int
	wclk    int
	is16    bool
	started bool
}

// NewSPI returns new SPI based implementation of tftdrv.DCI.
func NewSPI(drv *spi.Master, csn, dc iomux.Pin, mode spi.Config, rcHz, wcHz int) *SPI {
	dci := &SPI{
		spi:    drv,
		dc:     gpio.UsePin(dc),
		csn:    gpio.UsePin(csn),
		mode8:  mode | spi.Word8b,
		mode16: mode | spi.Word16b,
		rclk:   rcHz,
		wclk:   wcHz,
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

func set8b(dci *SPI) {
	if dci.is16 {
		dci.is16 = false
		dci.spi.SetConfig(dci.mode8)
	}
}

func set16b(dci *SPI) {
	if !dci.is16 {
		dci.is16 = true
		dci.spi.SetConfig(dci.mode16)
	}
}

func (dci *SPI) Cmd(p []byte, _ int) {
	if !dci.started {
		start(dci)
	}
	d := dci.spi
	d.WaitTxDone()
	dci.dc.Clear()
	set8b(dci)
	d.Write(p)
	d.WaitTxDone()
	dci.dc.Set()
}

// End ends the SPI transaction. It sets CSn pin high, disables the SPI
// peripheral and unlocks the driver. Other usesrs of the same master driver
// can then take controll of the SPI bus locking the driver before use.
func (dci *SPI) End() {
	if dci.started {
		dci.started = false
		dci.spi.Disable()
		dci.csn.Set()
		dci.spi.Unlock()
	}
}

func (dci *SPI) WriteBytes(p []uint8) {
	if !dci.started {
		start(dci)
	}
	set8b(dci)
	dci.spi.Write(p)
}

func (dci *SPI) WriteString(s string) {
	if !dci.started {
		start(dci)
	}
	set8b(dci)
	dci.spi.WriteString(s)
}

func (dci *SPI) WriteByteN(b byte, n int) {
	if !dci.started {
		start(dci)
	}
	set8b(dci)
	dci.spi.WriteByteN(b, n)
}

func (dci *SPI) WriteWords(p []uint16) {
	if !dci.started {
		start(dci)
	}
	set16b(dci)
	dci.spi.Write16(p)
}

func (dci *SPI) WriteWordN(w uint16, n int) {
	if !dci.started {
		start(dci)
	}
	set16b(dci)
	dci.spi.WriteWord16N(w, n)
}

func (dci *SPI) ReadBytes(p []byte) {
	if !dci.started {
		start(dci)
	}
	set8b(dci)
	d := dci.spi
	d.SetBaudrate(dci.rclk)
	d.Read(p)
	d.SetBaudrate(dci.wclk)
}
