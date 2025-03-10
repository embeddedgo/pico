// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package init when imported configures the whole system for typical usage
// assuming it's RPI Pico 2 compatible, that is, it's clocked from 12 MHz
// crystal, the XIP QSPI Flash supports 133 MHz clock, the IOVDD >= 2.5 V. Both
// the CPU and Flash are configured to run at conservative 125 MHz. The RISCV
// platform timer is used as the system timer.
//
// The CPU and Flash clock can be customized by using build tags. The 133MHz
// build tag configures both the CPU and Flash to run at 133 MHz. The 150MHz tag
// sets the CPU clock to 150 MHz and the Flash clock to 75 MHz.
//
// All peripheral drivers from the hal directory import this package to ensure
// proper system initialization. You can avoid any effects of importing this
// package by setting the nosysinit build flag.
package init
