// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pwr

import "github.com/embeddedgo/pico/hal/iomux"

const (
	VBUS = iomux.P24 // VBUS presence
	VSYS = iomux.P29 // analog VSYS/3
)

// SetPowerSave controls the onboard buck-boost DCD regulator power save mode.
// The default mode is ps=true which means the regulator is in Pulse Frequency
// Modulation at light loads and switches to the Pulse Width Modulation (PWM)
// mode only under heavy load. Setting ps to false forces PWM mode which reduces
// the output ripple (which may for example improve the ADC performance) at
// light load but at the expense of much worse efficiency.
func SetPowerSave(ps bool) {
	pin := iomux.P23
	pin.Setup(iomux.D4mA)
	if ps {
		pin.SetAltFunc(iomux.OutLow | iomux.OEEnable)
	} else {
		pin.SetAltFunc(iomux.OutHigh | iomux.OEEnable)
	}
}
