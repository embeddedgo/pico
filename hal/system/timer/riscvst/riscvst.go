// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riscvst

import (
	"embedded/rtos"
	"runtime"

	"github.com/embeddedgo/pico/hal/irq"
	"github.com/embeddedgo/pico/hal/system"
	"github.com/embeddedgo/pico/p/sio"
	"github.com/embeddedgo/pico/p/ticks"
)

// Setup configures and sets the RISCV platform timer as the tickless system
// timer. The timer resolution is 1 uS (true for any integer MHz crystal). .
func Setup() {
	runtime.LockOSThread()
	pl, _ := rtos.SetPrivLevel(0)

	t := &ticks.TICKS().T[ticks.RISCV]
	t.CYCLES.Store(uint32(system.ClockREF.Freq()) / 1e6)
	t.CTRL.Store(ticks.ENABLE)

	rtos.SetPrivLevel(pl)
	runtime.UnlockOSThread()

	SIO := sio.SIO()
	SIO.MTIMECMPH.Store(0xffff_ffff)
	SIO.MTIME_CTRL.Store(sio.EN)
	irq.SIO_IRQ_MTIMECMP.Enable(rtos.IntPrioSysTimer, 0)

	rtos.SetSystemTimer(nanotime, setAlarm)
}

//go:nosplit
func nanotime() int64 {
	SIO := sio.SIO()
	ph := SIO.MTIMEH.Load()
	l := SIO.MTIME.Load()
	h := SIO.MTIMEH.Load()
	if h != ph {
		l = 0
	}
	return (int64(h)<<32 + int64(l)) * 1e3
}

//go:nosplit
func setAlarm(ns int64) {
	timecmp := uint64(1<<64 - 1)
	if ns >= 0 {
		timecmp = uint64(ns) / 1e3
	}
	h := uint32(timecmp) >> 32
	l := uint32(timecmp)
	SIO := sio.SIO()
	SIO.MTIMECMPH.Store(0xffff_ffff)
	SIO.MTIMECMP.Store(l)
	SIO.MTIMECMPH.Store(h)
}
