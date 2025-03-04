// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uartcon

import (
	"embedded/rtos"
	"os"
	"syscall"

	"github.com/embeddedgo/fs/termfs"
	"github.com/embeddedgo/pico/hal/iomux"
	"github.com/embeddedgo/pico/hal/uart"
)

var ud *uart.Driver

func write(_ int, p []byte) int {
	n, _ := ud.Write(p)
	return n
}

func panicErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Setup setpus an LPUART peripheral to work as the system console.
func Setup(d *uart.Driver, rx, tx iomux.Pin, conf uart.Config, baudrate int, name string) {
	// Setup and enable the LPUART driver.
	d.UsePin(tx, uart.TXD)
	d.UsePin(rx, uart.RXD)
	d.Setup(conf, baudrate)
	d.EnableTx()
	d.EnableRx()

	// Set a system writer for print, println, panic, etc.
	ud = d
	rtos.SetSystemWriter(write)

	// Setup a serial console (standard input and output).
	con := termfs.New(name, d, d)
	con.SetCharMap(termfs.InCRLF | termfs.OutLFCRLF)
	con.SetEcho(true)
	con.SetLineMode(true, 256)
	rtos.Mount(con, "/dev/console")
	var err error
	os.Stdin, err = os.OpenFile("/dev/console", syscall.O_RDONLY, 0)
	panicErr(err)
	os.Stdout, err = os.OpenFile("/dev/console", syscall.O_WRONLY, 0)
	panicErr(err)
	os.Stderr = os.Stdout
}

// SetupLight setpus an LPUART to work as the system console.
// It usese termfs.LightFS instead of termfs.FS.
func SetupLight(d *uart.Driver, rx, tx iomux.Pin, conf uart.Config, baudrate int, name string) {
	// Setup and enable the LPUART driver.
	d.UsePin(tx, uart.TXD)
	d.UsePin(rx, uart.RXD)
	d.Setup(conf, baudrate)
	d.EnableTx()
	d.EnableRx()

	// Set a system writer for print, println, panic, etc.
	ud = d
	rtos.SetSystemWriter(write)

	// Setup a serial console (standard input and output).
	con := termfs.NewLight(name, d, d)
	rtos.Mount(con, "/dev/console")
	var err error
	os.Stdin, err = os.OpenFile("/dev/console", syscall.O_RDONLY, 0)
	panicErr(err)
	os.Stdout, err = os.OpenFile("/dev/console", syscall.O_WRONLY, 0)
	panicErr(err)
	os.Stderr = os.Stdout
}
