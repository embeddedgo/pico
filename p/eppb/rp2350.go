// DO NOT EDIT THIS FILE. GENERATED BY svdxgen.

//go:build rp2350

// Package eppb provides access to the registers of the EPPB peripheral.
//
// Instances:
//
//	EPPB  EPPB_BASE  -  -  Cortex-M33 EPPB vendor register block for RP2350
//
// Registers:
//
//	0x000 32  NMI_MASK0  NMI mask for IRQs 0 through 31. This register is core-local, and is reset by a processor warm reset.
//	0x004 32  NMI_MASK1  NMI mask for IRQs 0 though 51. This register is core-local, and is reset by a processor warm reset.
//	0x008 32  SLEEPCTRL  Nonstandard sleep control register
//
// Import:
//
//	github.com/embeddedgo/pico/p/mmap
package eppb

const (
	LIGHT_SLEEP SLEEPCTRL = 0x01 << 0 //+ By default, any processor sleep will deassert the system-level clock request. Reenabling the clocks incurs 5 cycles of additional latency on wakeup. Setting LIGHT_SLEEP to 1 keeps the clock request asserted during a normal sleep (Arm SCR.SLEEPDEEP = 0), for faster wakeup. Processor deep sleep (Arm SCR.SLEEPDEEP = 1) is not affected, and will always deassert the system-level clock request.
	WICENREQ    SLEEPCTRL = 0x01 << 1 //+ Request that the next processor deep sleep is a WIC sleep. After setting this bit, before sleeping, poll WICENACK to ensure the processor interrupt controller has acknowledged the change.
	WICENACK    SLEEPCTRL = 0x01 << 2 //+ Status signal from the processor's interrupt controller. Changes to WICENREQ are eventually reflected in WICENACK.
)

const (
	LIGHT_SLEEPn = 0
	WICENREQn    = 1
	WICENACKn    = 2
)