// DO NOT EDIT THIS FILE. GENERATED BY xgen.

//go:build rp2350

package uart

import (
	"embedded/mmio"
	"unsafe"

	"github.com/embeddedgo/pico/p/mmap"
)

type Periph struct {
	DR        mmio.R32[uint32]
	RSR       mmio.R32[RSR]
	_         [4]uint32
	FR        mmio.R32[FR]
	_         uint32
	ILPR      mmio.R32[uint32]
	IBRD      mmio.R32[uint32]
	FBRD      mmio.R32[uint32]
	LCR_H     mmio.R32[LCR_H]
	CR        mmio.R32[CR]
	IFLS      mmio.R32[IFLS]
	IMSC      mmio.R32[INT]
	RIS       mmio.R32[INT]
	MIS       mmio.R32[INT]
	ICR       mmio.R32[INT]
	DMACR     mmio.R32[DMACR]
	_         [997]uint32
	PERIPHID0 mmio.R32[uint32]
	PERIPHID1 mmio.R32[PERIPHID1]
	PERIPHID2 mmio.R32[PERIPHID2]
	PERIPHID3 mmio.R32[uint32]
	PCELLID0  mmio.R32[uint32]
	PCELLID1  mmio.R32[uint32]
	PCELLID2  mmio.R32[uint32]
	PCELLID3  mmio.R32[uint32]
}

func UART0() *Periph { return (*Periph)(unsafe.Pointer(uintptr(mmap.UART0_BASE))) }
func UART1() *Periph { return (*Periph)(unsafe.Pointer(uintptr(mmap.UART1_BASE))) }

func (p *Periph) BaseAddr() uintptr {
	return uintptr(unsafe.Pointer(p))
}

type RSR uint32

func FE_(p *Periph) mmio.RM32[RSR] { return mmio.RM32[RSR]{&p.RSR, FE} }
func PE_(p *Periph) mmio.RM32[RSR] { return mmio.RM32[RSR]{&p.RSR, PE} }
func BE_(p *Periph) mmio.RM32[RSR] { return mmio.RM32[RSR]{&p.RSR, BE} }
func OE_(p *Periph) mmio.RM32[RSR] { return mmio.RM32[RSR]{&p.RSR, OE} }

type FR uint32

func CTS_(p *Periph) mmio.RM32[FR]  { return mmio.RM32[FR]{&p.FR, CTS} }
func DSR_(p *Periph) mmio.RM32[FR]  { return mmio.RM32[FR]{&p.FR, DSR} }
func DCD_(p *Periph) mmio.RM32[FR]  { return mmio.RM32[FR]{&p.FR, DCD} }
func BUSY_(p *Periph) mmio.RM32[FR] { return mmio.RM32[FR]{&p.FR, BUSY} }
func RXFE_(p *Periph) mmio.RM32[FR] { return mmio.RM32[FR]{&p.FR, RXFE} }
func TXFF_(p *Periph) mmio.RM32[FR] { return mmio.RM32[FR]{&p.FR, TXFF} }
func RXFF_(p *Periph) mmio.RM32[FR] { return mmio.RM32[FR]{&p.FR, RXFF} }
func TXFE_(p *Periph) mmio.RM32[FR] { return mmio.RM32[FR]{&p.FR, TXFE} }
func RI_(p *Periph) mmio.RM32[FR]   { return mmio.RM32[FR]{&p.FR, RI} }

type LCR_H uint32

func BRK_(p *Periph) mmio.RM32[LCR_H]  { return mmio.RM32[LCR_H]{&p.LCR_H, BRK} }
func PEN_(p *Periph) mmio.RM32[LCR_H]  { return mmio.RM32[LCR_H]{&p.LCR_H, PEN} }
func EPS_(p *Periph) mmio.RM32[LCR_H]  { return mmio.RM32[LCR_H]{&p.LCR_H, EPS} }
func STP2_(p *Periph) mmio.RM32[LCR_H] { return mmio.RM32[LCR_H]{&p.LCR_H, STP2} }
func FEN_(p *Periph) mmio.RM32[LCR_H]  { return mmio.RM32[LCR_H]{&p.LCR_H, FEN} }
func WLEN_(p *Periph) mmio.RM32[LCR_H] { return mmio.RM32[LCR_H]{&p.LCR_H, WLEN} }
func SPS_(p *Periph) mmio.RM32[LCR_H]  { return mmio.RM32[LCR_H]{&p.LCR_H, SPS} }

type CR uint32

func UARTEN_(p *Periph) mmio.RM32[CR] { return mmio.RM32[CR]{&p.CR, UARTEN} }
func SIREN_(p *Periph) mmio.RM32[CR]  { return mmio.RM32[CR]{&p.CR, SIREN} }
func SIRLP_(p *Periph) mmio.RM32[CR]  { return mmio.RM32[CR]{&p.CR, SIRLP} }
func LBE_(p *Periph) mmio.RM32[CR]    { return mmio.RM32[CR]{&p.CR, LBE} }
func TXE_(p *Periph) mmio.RM32[CR]    { return mmio.RM32[CR]{&p.CR, TXE} }
func RXE_(p *Periph) mmio.RM32[CR]    { return mmio.RM32[CR]{&p.CR, RXE} }
func DTR_(p *Periph) mmio.RM32[CR]    { return mmio.RM32[CR]{&p.CR, DTR} }
func RTS_(p *Periph) mmio.RM32[CR]    { return mmio.RM32[CR]{&p.CR, RTS} }
func OUT1_(p *Periph) mmio.RM32[CR]   { return mmio.RM32[CR]{&p.CR, OUT1} }
func OUT2_(p *Periph) mmio.RM32[CR]   { return mmio.RM32[CR]{&p.CR, OUT2} }
func RTSEN_(p *Periph) mmio.RM32[CR]  { return mmio.RM32[CR]{&p.CR, RTSEN} }
func CTSEN_(p *Periph) mmio.RM32[CR]  { return mmio.RM32[CR]{&p.CR, CTSEN} }

type IFLS uint32

func TXIFLSEL_(p *Periph) mmio.RM32[IFLS] { return mmio.RM32[IFLS]{&p.IFLS, TXIFLSEL} }
func RXIFLSEL_(p *Periph) mmio.RM32[IFLS] { return mmio.RM32[IFLS]{&p.IFLS, RXIFLSEL} }

type INT uint32

type DMACR uint32

func RXDMAE_(p *Periph) mmio.RM32[DMACR]   { return mmio.RM32[DMACR]{&p.DMACR, RXDMAE} }
func TXDMAE_(p *Periph) mmio.RM32[DMACR]   { return mmio.RM32[DMACR]{&p.DMACR, TXDMAE} }
func DMAONERR_(p *Periph) mmio.RM32[DMACR] { return mmio.RM32[DMACR]{&p.DMACR, DMAONERR} }

type PERIPHID1 uint32

func PARTNUMBER1_(p *Periph) mmio.RM32[PERIPHID1] {
	return mmio.RM32[PERIPHID1]{&p.PERIPHID1, PARTNUMBER1}
}
func DESIGNER0_(p *Periph) mmio.RM32[PERIPHID1] { return mmio.RM32[PERIPHID1]{&p.PERIPHID1, DESIGNER0} }

type PERIPHID2 uint32

func DESIGNER1_(p *Periph) mmio.RM32[PERIPHID2] { return mmio.RM32[PERIPHID2]{&p.PERIPHID2, DESIGNER1} }
func REVISION_(p *Periph) mmio.RM32[PERIPHID2]  { return mmio.RM32[PERIPHID2]{&p.PERIPHID2, REVISION} }
