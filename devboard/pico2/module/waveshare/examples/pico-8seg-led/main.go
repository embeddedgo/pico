// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/module/waveshare/pico-8seg-led/segled"
)

func main() {
	d := segled.Display
	for i := -99; ; i++ {
		d.Clear()
		fmt.Fprintf(d, "%4d\n", i)
		time.Sleep(time.Second / 4)
	}
}
