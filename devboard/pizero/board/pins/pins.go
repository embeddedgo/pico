// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pins

import "github.com/embeddedgo/pico/hal/iomux"

// Board markings
const (
	GP0  = iomux.P00
	GP1  = iomux.P01
	GP2  = iomux.P02
	GP3  = iomux.P03
	GP4  = iomux.P04
	GP5  = iomux.P05
	GP6  = iomux.P06
	GP7  = iomux.P07
	GP8  = iomux.P08
	GP9  = iomux.P09
	GP10 = iomux.P10
	GP11 = iomux.P11
	GP12 = iomux.P12
	GP13 = iomux.P13
	GP14 = iomux.P14
	GP15 = iomux.P15
	GP16 = iomux.P16
	GP17 = iomux.P17
	GP18 = iomux.P18
	GP19 = iomux.P19
	GP20 = iomux.P20
	GP21 = iomux.P21
	GP22 = iomux.P22
	GP23 = iomux.P23
	GP24 = iomux.P24
	GP25 = iomux.P25
	GP26 = iomux.P26
	GP27 = iomux.P27
)

// Broadcom (BCM) numbering
const (
	BCM2  = GP2
	BCM3  = GP3
	BCM4  = GP14 // differs
	BCM17 = GP17
	BCM27 = GP27
	BCM22 = GP22
	BCM10 = GP11 // differs
	BCM9  = GP12 // differs
	BCM11 = GP10 // differs
	BCM0  = GP0
	BCM5  = GP15 // differs
	BCM6  = GP6
	BCM13 = GP13
	BCM19 = GP19
	BCM26 = GP26

	BCM14 = GP4 // differs
	BCM15 = GP5 // differs
	BCM18 = GP18
	BCM23 = GP23
	BCM24 = GP24
	BCM25 = GP25
	BCM8  = GP8
	BCM7  = GP7
	BCM1  = GP1
	BCM12 = GP9 // differs
	BCM16 = GP16
	BCM20 = GP20
	BCM21 = GP21
)

// Wiring Pi numbering
const (
	WPi8  = GP2
	WPi9  = GP3
	WPi7  = GP14
	WPi0  = GP17
	WPi2  = GP27
	WPi3  = GP22
	WPi12 = GP11
	WPi13 = GP12
	WPi14 = GP10
	// GP0 has no WPi mapping
	WPi21 = GP15
	WPi22 = GP6
	WPi23 = GP13
	WPi24 = GP19
	WPi25 = GP26

	WPi15 = GP4
	WPi16 = GP5
	WPi1  = GP18
	WPi4  = GP23
	WPi5  = GP24
	// GP1 has no WPi mapping
	WPi26 = GP9
	WPi27 = GP16
	WPi28 = GP20
	WPi29 = GP21
)
