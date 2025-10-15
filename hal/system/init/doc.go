// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package init when imported configures the whole system for typical usage
// assuming it's RPI Pico 2 compatible, that is, it's clocked from 12 MHz
// crystal, the XIP QSPI Flash supports 133 MHz clock, the IOVDD >= 2.5 V. Both
// the CPU and flash are configured to run at the 125 MHz.. The RISCV platform
// timer is used as the system timer.
//
// The CPU and Flash clock can be customized by using build tags. The 133MHz
// build tag configures both the CPU and Flash to run at 133 MHz which may
// expose problems with some cheap Pico 2 clones (the reason is mainly their
// flash memory). The 150MHz tag sets the CPU clock to 150 MHz and the flash
// clock to very forgiving 75 MHz. The 200MHz tag overclocks the CPU and most
// peripherals to 200 MHz, the flash clock is 100 MHz.
//
// The flash clock is also used for the PSRAM if your board has one. Keep in
// mind that most common PSRAM chips support up 100 MHz clock.
//
// All peripheral drivers from the hal directory import this package to ensure
// proper system initialization. You can avoid any effects of importing this
// package by setting the nosysinit build flag.
package init
