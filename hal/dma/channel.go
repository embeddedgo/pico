// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dma

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/hal/internal"
)

// Channel rperesents a DMA channel.
type Channel struct {
	d *Controller
	n int
}

// IsValid reports whether the channel is valid, e.g. succesfully allocated by
// the Controller.AllocChannel method.
func (c Channel) IsValid() bool {
	return c.d != nil
}

// Controller returns the controller that the channel belongs to.
func (c Channel) Controller() *Controller {
	return c.d
}

// Num return the channel number in the controller.
func (c Channel) Num() int {
	return c.n
}

// Free frees the channel so the Controller.AllocChannel can allocate it next
// time.
func (c Channel) Free() {
	mask := uint32(1) << uint(c.n)
	chanAlloc.mx.Lock()
	chanAlloc.mask |= mask
	chanAlloc.mx.Unlock()
}

func addr(r *mmio.U32) uintptr {
	return uintptr(r.Load())
}

func setAddr(r *mmio.U32, a unsafe.Pointer) {
	r.Store(uint32(uintptr(a)))
}

func transCount(r *mmio.U32) (cnt int, mode int8) {
	v := r.Load()
	return int(v & 0x0fff_ffff), int8(v >> 28)
}

func setTransCount(r *mmio.U32, cnt int, mode int8) {
	r.Store(uint32(cnt)&0x0fff_ffff | uint32(mode)<<28)
}

func conf(r *mmio.U32, c Channel) (cfg Config, chainTo Channel) {
	v := r.Load()
	cfg = Config(v & 0x03fe_1fff)
	n := int(v>>13) & 15
	if n != c.n {
		chainTo.d = c.d
		chainTo.n = n
	}
	return
}

func setConf(r *mmio.U32, c Channel, cfg Config, chainTo Channel) {
	n := c.n // by default chain to itself (chaining disabled)
	if chainTo.d != nil {
		if chainTo.d != c.d {
			panic("dma: chainTo")
		}
		n = chainTo.n
	}
	r.Store(uint32(cfg) | uint32(n)<<13)
}

func status(r *mmio.U32) uint8 {
	return uint8(r.Load() >> 25)
}

func clear(r *mmio.U32, status uint8) {
	internal.AtomicSetU32(r, uint32(status)<<25)
}

func (c Channel) ReadAddr() uintptr {
	return addr(&c.d.ch[c.n].readAddr)
}

func (c Channel) SetReadAddr(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].readAddr, a)
}

func (c Channel) SetReadAddrTrig(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].readAddrTrig3, a)
}

func (c Channel) WriteAddr() uintptr {
	return addr(&c.d.ch[c.n].writeAddr)
}

func (c Channel) SetWriteAddr(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].writeAddr, a)
}

func (c Channel) SetWriteAddrTrig(a unsafe.Pointer) {
	setAddr(&c.d.ch[c.n].writeAddrTrig2, a)
}

// Transfer counter mode
const (
	Normal      int8 = 0  // normal mode, compatible with the RP2040
	TriggerSelf int8 = 1  // re-triggers itself at the end of transfer
	Endless     int8 = 15 // perform an endless sequence of transfers
)

// TransCount return the current transfer counter and mode.
func (c Channel) TransCount() (cnt int, mode int8) {
	return transCount(&c.d.ch[c.n].transCount)
}

// SetTransCount sets the transfer counter and mode.
func (c Channel) SetTransCount(cnt int, mode int8) {
	setTransCount(&c.d.ch[c.n].transCount, cnt, mode)
}

// SetTransCountTrig is like SetTransCount but also triggers the transfer.
func (c Channel) SetTransCountTrig(cnt int, mode int8) {
	setTransCount(&c.d.ch[c.n].transCountTrig1, cnt, mode)
}

// Config represents the channel configuration.
type Config uint32

const (
	En    Config = 1 << 0 // DMA channel enable
	PrioH Config = 1 << 1 // Referential treatment in issue scheduling

	DataSize Config = 3 << 2 // Size of each bus transfer:
	S8b      Config = 0 << 2 // - byte
	S16b     Config = 1 << 2 // - half word
	S32b     Config = 2 << 2 // - word

	IncR  Config = 1 << 4  // Increment read address with each transfer
	RevR  Config = 1 << 5  // Decrement read address rather than increment
	IncW  Config = 1 << 6  // Increment write address with each transfer
	RevW  Config = 1 << 7  // Decrement write address rather than increment
	RingW Config = 1 << 12 // Apply RingSize to write instead of read

	TransReq   Config = 0x3F << 17 // Select a Transfer Request signal:
	PIO0_TX0   Config = 0x00 << 17 // - PIO0's TX FIFO 0
	PIO0_TX1   Config = 0x01 << 17 // - PIO0's TX FIFO 1
	PIO0_TX2   Config = 0x02 << 17 // - PIO0's TX FIFO 2
	PIO0_TX3   Config = 0x03 << 17 // - PIO0's TX FIFO 3
	PIO0_RX0   Config = 0x04 << 17 // - PIO0's RX FIFO 0
	PIO0_RX1   Config = 0x05 << 17 // - PIO0's RX FIFO 1
	PIO0_RX2   Config = 0x06 << 17 // - PIO0's RX FIFO 2
	PIO0_RX3   Config = 0x07 << 17 // - PIO0's RX FIFO 3
	PIO1_TX0   Config = 0x08 << 17 // - PIO1's TX FIFO 0
	PIO1_TX1   Config = 0x09 << 17 // - PIO1's TX FIFO 1
	PIO1_TX2   Config = 0x0A << 17 // - PIO1's TX FIFO 2
	PIO1_TX3   Config = 0x0B << 17 // - PIO1's TX FIFO 3
	PIO1_RX0   Config = 0x0C << 17 // - PIO1's RX FIFO 0
	PIO1_RX1   Config = 0x0D << 17 // - PIO1's RX FIFO 1
	PIO1_RX2   Config = 0x0E << 17 // - PIO1's RX FIFO 2
	PIO1_RX3   Config = 0x0F << 17 // - PIO1's RX FIFO 3
	PIO2_TX0   Config = 0x10 << 17 // - PIO2's TX FIFO 0
	PIO2_TX1   Config = 0x11 << 17 // - PIO2's TX FIFO 1
	PIO2_TX2   Config = 0x12 << 17 // - PIO2's TX FIFO 2
	PIO2_TX3   Config = 0x13 << 17 // - PIO2's TX FIFO 3
	PIO2_RX0   Config = 0x14 << 17 // - PIO2's RX FIFO 0
	PIO2_RX1   Config = 0x15 << 17 // - PIO2's RX FIFO 1
	PIO2_RX2   Config = 0x16 << 17 // - PIO2's RX FIFO 2
	PIO2_RX3   Config = 0x17 << 17 // - PIO2's RX FIFO 3
	SPI0_TX    Config = 0x18 << 17 // - SPI0's TX FIFO
	SPI0_RX    Config = 0x19 << 17 // - SPI0's RX FIFO
	SPI1_TX    Config = 0x1A << 17 // - SPI1's TX FIFO
	SPI1_RX    Config = 0x1B << 17 // - SPI1's RX FIFO
	UART0_TX   Config = 0x1C << 17 // - UART0's TX FIFO
	UART0_RX   Config = 0x1D << 17 // - UART0's RX FIFO
	UART1_TX   Config = 0x1E << 17 // - UART1's TX FIFO
	UART1_RX   Config = 0x1F << 17 // - UART1's RX FIFO
	PWM_WRAP0  Config = 0x20 << 17 // - PWM Counter 0's Wrap Value
	PWM_WRAP1  Config = 0x21 << 17 // - PWM Counter 1's Wrap Value
	PWM_WRAP2  Config = 0x22 << 17 // - PWM Counter 2's Wrap Value
	PWM_WRAP3  Config = 0x23 << 17 // - PWM Counter 3's Wrap Value
	PWM_WRAP4  Config = 0x24 << 17 // - PWM Counter 4's Wrap Value
	PWM_WRAP5  Config = 0x25 << 17 // - PWM Counter 5's Wrap Value
	PWM_WRAP6  Config = 0x26 << 17 // - PWM Counter 6's Wrap Value
	PWM_WRAP7  Config = 0x27 << 17 // - PWM Counter 7's Wrap Value
	PWM_WRAP8  Config = 0x28 << 17 // - PWM Counter 8's Wrap Value
	PWM_WRAP9  Config = 0x29 << 17 // - PWM Counter 9's Wrap Value
	PWM_WRAP10 Config = 0x2A << 17 // - PWM Counter 0's Wrap Value
	PWM_WRAP11 Config = 0x2B << 17 // - PWM Counter 1's Wrap Value
	I2C0_TX    Config = 0x2C << 17 // - I2C0's TX FIFO
	I2C0_RX    Config = 0x2D << 17 // - I2C0's RX FIFO
	I2C1_TX    Config = 0x2E << 17 // - I2C1's TX FIFO
	I2C1_RX    Config = 0x2F << 17 // - I2C1's RX FIFO
	ADC        Config = 0x30 << 17 // - the ADC
	XIP_STREAM Config = 0x31 << 17 // - the XIP Streaming FIFO
	XIP_QMITX  Config = 0x32 << 17 // - XIP_QMITX
	XIP_QMIRX  Config = 0x33 << 17 // - XIP_QMIRX
	HSTX       Config = 0x34 << 17 // - HSTX
	CORESIGHT  Config = 0x35 << 17 // - CORESIGHT
	SHA256     Config = 0x36 << 17 // - SHA256
	TIMER0     Config = 0x3B << 17 // - Timer 0
	TIMER1     Config = 0x3C << 17 // - Timer 1
	TIMER2     Config = 0x3D << 17 // - Timer 2 (Optional)
	TIMER3     Config = 0x3E << 17 // - Timer 3 (Optional)
	Always     Config = 0x3F << 17 // - permanent request

	Quiet  Config = 1 << 23 // IRQ only when zero/null written to trig register
	Swap   Config = 1 << 24 // Reverse the order of bytes in the transfered word
	SnifEn Config = 1 << 25 // Transfers are visible to the sniff hardware
)

// RingSize allows to spefiy the ring size parameter in the configuration.
func RingSize(log2size int) (c Config) {
	if uint(log2size) > 15 {
		panic("dma: log2size")
	}
	return Config(log2size) << 8
}

// RingSize returns the value of ring size parameter in the configuration.
func (cfg Config) RingSize() (log2size int) {
	return int(cfg >> 8 & 15)
}

// Conf returns the channel current configuration.
func (c Channel) Config() (cfg Config, chainTo Channel) {
	return conf(&c.d.ch[c.n].ctrlTrig, c)
}

// SetConfig configures the channel c according to the cfg with the chainTo
// channel to be triggered at the end of transfer. Set chainTo to Channel{} or
// c (chain to itself) to disable chaining.
func (c Channel) SetConfig(cfg Config, chainTo Channel) {
	setConf(&c.d.ch[c.n].ctrl1, c, cfg, chainTo)
}

// SetConfigTrig is like SetConfig but also triggers the transfer.
func (c Channel) SetConfigTrig(cfg Config, chainTo Channel) {
	setConf(&c.d.ch[c.n].ctrlTrig, c, cfg, chainTo)
}

// Status
const (
	Busy     uint8 = 1 << 0 // The channel performs transfer
	WriteErr uint8 = 1 << 3 // Write bus error (use Clear to clear)
	ReadErr  uint8 = 1 << 4 // Read bus error (
	AHBErr   uint8 = 1 << 5 // Logical OR of the WriteErr and ReadErr flags
)

// Status returns the channel status flags.
func (c Channel) Status() uint8 {
	return status(&c.d.ch[c.n].ctrlTrig)
}

// Clear clears the selected status bits. Only WriteErr and ReadErr bits an be
// cleard this way.
func (c Channel) Clear(status uint8) {
	clear(&c.d.ch[c.n].ctrl1, status)
}

// ClearTrig is like Clear but also triggers the transfer.
func (c Channel) ClearTrig(status uint8) {
	clear(&c.d.ch[c.n].ctrlTrig, status)
}

// Trig triggers the transfer.
func (c Channel) Trig() {
	c.d.multiChanTrig.Store(1 << uint(c.n))
}

// EnableIRQ enables unsing interrup line irqn to sending interrupt request.
func (c Channel) EnableIRQ(irqn int) {
	if uint(irqn) > 3 {
		panic("dma: irqn")
	}
	internal.AtomicSetU32(&c.d.irq[irqn].e, 1<<uint(c.n))
}

// DisableIRQ disables unsing interrup line irqn to sending interrupt request.
func (c Channel) DisableIRQ(irqn int) {
	if uint(irqn) > 3 {
		panic("dma: irqn")
	}
	internal.AtomicClearU32(&c.d.irq[irqn].e, 1<<uint(c.n))
}

// IRQEnabled reports whether the interrupt line irqn is used by this channel.
func (c Channel) IRQEnabled(irqn int) bool {
	if uint(irqn) > 3 {
		panic("dma: irqn")
	}
	return c.d.irq[irqn].e.LoadBits(1<<uint(c.n)) != 0
}

// IsIRQ reports whether the interrupt request is active for this channel.
func (c Channel) IsIRQ() bool {
	return c.d.irq[0].r.LoadBits(1<<uint(c.n)) != 0
}

// ClearIRQ clears the interrupr request for this channel.
func (c Channel) ClearIRQ() {
	c.d.irq[0].r.Store(1 << uint(c.n))
}
