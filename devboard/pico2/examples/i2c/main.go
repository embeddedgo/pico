package main

import (
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/leds"
	"github.com/embeddedgo/pico/hal/i2c"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/system/clock"
)

func setBaudrate(p *i2c.Periph, baud int) {
	clk := clock.PERI.Freq()
	cn := uint32((clk + int64(baud/2)) / int64(baud))

	//spkLen := p.FS_SPKLEN.Load()
	//lcn := cn/2 - 1
	//hcn := cn/2  - spkLen - 7
	lcn := cn * 3 / 5
	hcn := cn - lcn

	var txHold i2c.SDA_HOLD
	if baud < 1e6 {
		txHold = i2c.SDA_HOLD(clk*3/10e6) + 1
	} else {
		txHold = i2c.SDA_HOLD(clk*3/25e6) + 1
	}

	spkn := uint32(1)
	if lcn >= 16 {
		spkn = lcn / 16
	}

	p.FS_SCL_LCNT.Store(lcn)
	p.FS_SCL_HCNT.Store(hcn)
	p.FS_SPKLEN.Store(spkn)
	p.SDA_HOLD.StoreBits(i2c.SDA_TX_HOLD, txHold)
}

func main() {
	sda := iomux.P20
	scl := iomux.P21

	sda.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	sda.SetAltFunc(iomux.I2C)
	scl.Setup(iomux.InpEn | iomux.D4mA | iomux.PullUp)
	scl.SetAltFunc(iomux.I2C)

	p := i2c.I2C(0)
	p.SetReset(true)
	p.SetReset(false)

	p.ENABLE.Store(0)

	p.CON.Store(i2c.FAST | i2c.MASTER_MODE | i2c.SLAVE_DISABLE | i2c.RESTART_EN | i2c.TX_EMPTY_CTRL)

	p.TX_TL.Store(0)
	p.RX_TL.Store(0)

	setBaudrate(p, 100e3)

	p.TAR.Store(0b010_0111)
	p.ENABLE.Store(i2c.EN)

	for b := uint8(0); ; b++ {
		for p.STATUS.LoadBits(i2c.TFNF) == 0 {
		}
		p.DATA_CMD.Store(i2c.STOP | uint32(b))
		if true || b == 0 {
			leds.User.Toggle()
		}
		time.Sleep(time.Second / 2)
	}
}
