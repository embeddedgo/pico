// DO NOT EDIT THIS FILE. GENERATED BY xgen.

//go:build rp2350

package padsbank

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/p/mmap"
)

type Periph struct {
	VOLTAGE_SELECT mmio.R32[VOLTAGE_SELECT]
	GPIO           [48]mmio.R32[GPIO]
	SWCLK          mmio.R32[GPIO]
	SWD            mmio.R32[GPIO]
}

func PADS_BANK0() *Periph { return (*Periph)(unsafe.Pointer(uintptr(mmap.PADS_BANK0_BASE))) }

func (p *Periph) BaseAddr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

type VOLTAGE_SELECT uint32

func IOVDD_(p *Periph) mmio.RM32[VOLTAGE_SELECT] {
	return mmio.RM32[VOLTAGE_SELECT]{&p.VOLTAGE_SELECT, IOVDD}
}

type GPIO uint32

func SLEWFAST_(p *Periph, i int) mmio.RM32[GPIO] { return mmio.RM32[GPIO]{&p.GPIO[i], SLEWFAST} }
func SCHMITT_(p *Periph, i int) mmio.RM32[GPIO]  { return mmio.RM32[GPIO]{&p.GPIO[i], SCHMITT} }
func PDE_(p *Periph, i int) mmio.RM32[GPIO]      { return mmio.RM32[GPIO]{&p.GPIO[i], PDE} }
func PUE_(p *Periph, i int) mmio.RM32[GPIO]      { return mmio.RM32[GPIO]{&p.GPIO[i], PUE} }
func DRIVE_(p *Periph, i int) mmio.RM32[GPIO]    { return mmio.RM32[GPIO]{&p.GPIO[i], DRIVE} }
func IE_(p *Periph, i int) mmio.RM32[GPIO]       { return mmio.RM32[GPIO]{&p.GPIO[i], IE} }
func OD_(p *Periph, i int) mmio.RM32[GPIO]       { return mmio.RM32[GPIO]{&p.GPIO[i], OD} }
func ISO_(p *Periph, i int) mmio.RM32[GPIO]      { return mmio.RM32[GPIO]{&p.GPIO[i], ISO} }
