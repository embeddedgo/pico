// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the code from the pico-sdk (Copyright (c) 2020 Raspberry
// Pi (Trading) Ltd; SPDX-License-Identifier: BSD-3-Clause) loosely translated
// to Go.

package system

import (
	"embedded/rtos"
	"runtime"

	"github.com/embeddedgo/pico/hal/internal"
	"github.com/embeddedgo/pico/p/clocks"
	"github.com/embeddedgo/pico/p/pll"
	"github.com/embeddedgo/pico/p/qmi"
	"github.com/embeddedgo/pico/p/resets"
	"github.com/embeddedgo/pico/p/xosc"
)

// SetupPico2_125MHz setups the system assuming it is an RPI Pico 2 compatible,
// that is, it's clocked from 12 MHz crystal, the XIP QSPI flash supports 133
// MHz clock, the IOVDD >= 2.5 V. Both the CPU and QSPI flash are configured
// to run at conservative 125 MHz, which seems to be a good compromise between
// speed and stability. Works well with Pico 2 clones that may have worse flash
// chip/circuit. See the Setup function for more generic interface.
func SetupPico2_125MHz() {
	Setup(12e6, PLL{1, 125, 6, 2}, PLL{1, 100, 5, 5}, 133e6)
}

// SetupPico2_133MHz is like SetupPico2_125MHz but both the CPU and QSPI flash
// run at 133 MHz. It gives you slightly faster system but without any  margin
// for flash and its circuity which may cause problems in case of cheap Pico 2
// clones.
func SetupPico2_133MHz() {
	Setup(12e6, PLL{1, 133, 6, 2}, PLL{1, 100, 5, 5}, 133e6)
}

// SetupPico2_150MHz is like SetupPico2_125MHz but the CPU runs at 150 MHz, the
// QMI (flash, PSRAM) at 75 MHz. It should give you quite good performance while
// being very forgiving for the flash and its circuity. It also supports most
// common QSPI PSRAM chips thet usually can work up to 100 MHz.
func SetupPico2_150MHz() {
	Setup(12e6, PLL{1, 125, 5, 2}, PLL{1, 100, 5, 5}, 133e6)
}

// SetupPico2_200MHz is like SetupPico2_150MHz but the CPU and most peripherals
// are overclocked to 200 MHz. The QMI (flash, PSRAM) run at 100 MHz which is
// below the spec of the Pico 2 flash but seems to be the maximum supported
// clock in case of most common PSRAM chips. This level of overclocking should
// work with almost all RP2350 chips but keep in mind that it's out of spec.
func SetupPico2_200MHz() {
	Setup(12e6, PLL{1, 100, 3, 2}, PLL{1, 100, 5, 5}, 133e6)
}

// SetupPico2_250MHz is like SetupPico2_125MHz but the CPU and most peripherals
// are overclocked to 250 MHz. The QMI (flash, PSRAM) runs at 125 MHz. It is
// significant level of overclocking but seems to be still well supported.
func SetupPico2_250MHz() {
	Setup(12e6, PLL{1, 125, 3, 2}, PLL{1, 100, 5, 5}, 133e6)
}

// A PLL configuration.
//
//	vcoHz = refHz / RefDiv * FbDiv
//	outHz = vcoHz / (PostDiv1 * PostDiv2)
//
// Constraints:
//
//	1 <= RefDiv <= 63
//	16 <= FbDiv <= 320
//	1 <= PostDiv <= 7
//	refHz / RefDiv >= 5 MHz
//	750 MHz <= vcoHz <= 1600 MHz
type PLL struct {
	RefDiv   int
	FbDiv    int
	PostDiv1 int
	PostDiv2 int
}

// Fout calculates the output frequency of the PLL worknig with the pll
// configuration and the refHz frequency as an input. It returns outHz < 0 if
// the refHz is invalid or the pll configuration is invalid for the given refHz.
func (pll PLL) Fout(refHz int64) (outHz int64) {
	if pll.RefDiv < 1 || 63 < pll.RefDiv {
		return -1
	}
	if pll.FbDiv < 16 || 320 < pll.FbDiv {
		return -1
	}
	if pll.PostDiv1 < 1 || 7 < pll.PostDiv1 {
		return -1
	}
	if pll.PostDiv2 < 1 || 7 < pll.PostDiv2 {
		return -1
	}
	if uint64(refHz) > 100e6 {
		return -1
	}
	frefHz := int(refHz) / pll.RefDiv
	if frefHz < 5e6 {
		return -1
	}
	vcoHz := frefHz * pll.FbDiv
	if vcoHz < 750e6 || 1600e6 < vcoHz {
		return -1
	}
	return int64(vcoHz / (pll.PostDiv1 * pll.PostDiv2))
}

/*
  pico-sdk initialization sequence

  pico_crt0/crt0.S   _reset_handler
  newlib_interface.c   runtime_init
  runtime.c              runtime_run_initializers
  runtime.c                runtime_run_initializers_from
  runtime_init.c             runtime_init_bootrom_reset
  runtime_init.c             runtime_init_early_resets
  runtime_init.c             runtime_init_usb_power_down
  runtime_init_clocks.c      runtime_init_clocks
  runtime_init.c             runtime_init_post_clock_resets
  boot_lock.c                runtime_init_boot_locks_reset
  runtime_init.c             runtime_init_spin_locks_reset
  bootrom_lock.c             runtime_init_bootrom_locking_enable
  mutex.c                    runtime_init_mutex
  runtime_init.c             runtime_init_install_ram_vector_table
  time.c                     runtime_init_default_alarm_pool
  runtime.c                  first_per_core_initializer
  runtime_init.c             runtime_init_per_core_bootrom_reset
  runtime_init.c             runtime_init_per_core_enable_coprocessors
  sync_spin_lock.c           spinlock_set_extexclall
  irq.c                      runtime_init_per_core_irq_priorities
*/

// Setup initializes and configures the system. It expects a crystal oscillator
// as the source of the reference frequency to both PLLs and an QSPI flash
// supporting maxFlashHz clock from which the code is executed (XIP).
func Setup(xoscHz int64, sys, usb PLL, maxFlashHz int64) {
	// The default configuration in the ACCESSCTRL makes some of the registers
	// used below only accessible in the priviledged mode.
	runtime.LockOSThread()
	pl, _ := rtos.SetPrivLevel(0)

	// 00050 PICO_RUNTIME_INIT_BOOTROM_RESET

	// TODO

	// 00100 PICO_RUNTIME_INIT_EARLY_RESETS

	// Reset all peripherals except these that may hurt the code execution.
	rst := resets.RESETS()
	allp := uint32(0x1fff_ffff) // all defined peripheral bits
	nrp := resets.IO_QSPI |
		resets.PADS_QSPI |
		resets.PLL_USB |
		resets.USBCTRL |
		resets.SYSCFG |
		resets.PLL_SYS
	internal.AtomicSet(&rst.RESET, allp&^nrp)

	// Remove reset from peripherals which are clocked only by SYS, REF.
	nup := resets.HSTX |
		resets.ADC |
		resets.SPI0 |
		resets.SPI1 |
		resets.UART0 |
		resets.UART1 |
		resets.USBCTRL
	internal.AtomicClear(&rst.RESET, allp&^nup)
	for rst.RESET_DONE.LoadBits(allp&^nup) != allp&^nup {
	}

	// TODO: runtime_init_usb_power_down

	// 00500 PICO_RUNTIME_INIT_CLOCKS

	// Disable resus that may be enabled from previous software.
	clk := clocks.CLOCKS()
	clk.SYS_RESUS_CTRL.Store(0)

	// Enable the xosc.
	var fr xosc.CTRL
	switch {
	case xoscHz < 15e6:
		fr = xosc.FR1_15MHZ
	case xoscHz < 30e6:
		fr = xosc.FR10_30MHZ
	case xoscHz < 60e6:
		fr = xosc.FR25_60MHZ
	default:
		fr = xosc.FR40_100MHZ
	}
	osc := xosc.XOSC()
	osc.CTRL.Store(fr)
	const delayMultipler = 64
	osc.STARTUP.Store(xosc.STARTUP((xoscHz/1e3 + 128) / 256 * delayMultipler))
	internal.AtomicSet(&osc.CTRL, xosc.ENABLE)
	for osc.STATUS.LoadBits(xosc.STABLE) == 0 {
	}

	// Switch SYS and REF cleanly away from their aux sources.
	internal.AtomicClear(&clk.CLK[clocks.SYS].CTRL, clocks.SYS_SRC)
	for clk.CLK[clocks.SYS].SELECTED.LoadBits(0b11) != 0b01 {
	}
	internal.AtomicClear(&clk.CLK[clocks.REF].CTRL, clocks.REF_SRC)
	for clk.CLK[clocks.REF].SELECTED.LoadBits(0b1111) != 0b0001 {
	}

	// Setup PLLs.
	sysHz := sys.Fout(xoscHz)
	if sysHz < 0 {
		panic("bad PLL_SYS cfg")
	}
	setupPLL(pll.SYS(), resets.PLL_SYS, sys)
	usbHz := usb.Fout(xoscHz)
	if usbHz < 0 {
		panic("bad PLL_USB cfg")
	}
	setupPLL(pll.USB(), resets.PLL_USB, usb)

	// Configure all 6 clocks.
	setClock(
		clocks.REF,
		clocks.REF_XOSC_CLKSRC, 0,
		uint(xoscHz), 1<<clocks.REF_INTn,
	)
	setClock(
		clocks.SYS,
		clocks.SYS_CLKSRC_CLK_SYS_AUX, clocks.SYS_CLKSRC_PLL_SYS,
		uint(sysHz), 1<<clocks.SYS_INTn,
	)
	setClock(
		clocks.PERI,
		0, clocks.PERI_CLK_SYS,
		uint(sysHz), 1<<clocks.PERI_INTn,
	)
	setClock(
		clocks.HSTX,
		0, clocks.HSTX_CLK_SYS,
		uint(sysHz), 1<<clocks.HSTX_INTn,
	)
	setClock(
		clocks.USB,
		0, clocks.USB_CLKSRC_PLL_USB,
		uint(usbHz), 1<<clocks.USB_INTn,
	)
	setClock(
		clocks.ADC,
		0, clocks.ADC_CLKSRC_PLL_USB,
		uint(usbHz), 1<<clocks.ADC_INTn,
	)

	// pico-sdk starts all tick generators here, all configurrd to 1 MHz. We
	// leave them disabled and enable one by one if needed.

	// 00600 PICO_RUNTIME_INIT_POST_CLOCK_RESETS

	// Remove reset from all peripherals
	internal.AtomicClear(&rst.RESET, allp)
	for rst.RESET_DONE.LoadBits(allp) != allp {
	}

	// Increase the QSPI Flash clock speed.
	qmiDiv := (uint(sysHz)-1)/uint(maxFlashHz) + 1
	qmi.QMI().M[0].TIMING.StoreBits(
		qmi.CLKDIV|qmi.RXDELAY,
		qmi.TIMING(qmiDiv)<<qmi.CLKDIVn| // SCK = freq(SYS) / CLKDIV
			qmi.TIMING(qmiDiv)<<qmi.RXDELAYn, // RXDELAY unit is period(SYS)/2
	)

	rtos.SetPrivLevel(pl)
	runtime.UnlockOSThread()
}

func setupPLL(p *pll.Periph, reset uint32, cfg PLL) {
	// Avoid disrupting PLL that is already correctly configured.
	pdiv := cfg.PostDiv1<<pll.POSTDIV1n | cfg.PostDiv2<<pll.POSTDIV2n
	cs := p.CS.Load()
	if cs&pll.LOCK != 0 &&
		cs&pll.REFDIV>>pll.REFDIVn == pll.CS(cfg.RefDiv) &&
		p.FBDIV_INT.Load() == uint32(cfg.FbDiv) &&
		p.PRIM.Load() == uint32(pdiv) {
		return
	}

	// Reset.
	rst := resets.RESETS()
	internal.AtomicSet(&rst.RESET, reset)
	internal.AtomicClear(&rst.RESET, reset)
	for rst.RESET_DONE.LoadBits(reset) == 0 {
	}

	// Turn on this PLL.
	p.CS.Store(pll.CS(cfg.RefDiv))
	p.FBDIV_INT.Store(uint32(cfg.FbDiv))
	internal.AtomicClear(&p.PWR, pll.PD|pll.VCOPD)
	for p.CS.LoadBits(pll.LOCK) == 0 {
	}

	// Turn on the post divider.
	p.PRIM.Store(uint32(pdiv))
	internal.AtomicClear(&p.PWR, pll.POSTDIVPD)
}

func setClock(clk int, src, auxsrc clocks.CTRL, freqHz uint, div clocks.DIV)
