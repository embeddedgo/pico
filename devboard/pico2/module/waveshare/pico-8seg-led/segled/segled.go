// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package segled

import (
	"time"

	"github.com/embeddedgo/pico/devboard/pico2/board/pins"
	"github.com/embeddedgo/pico/hal/gpio"
	"github.com/embeddedgo/pico/hal/iomux"
)

type ShiftReg struct {
	din, clk, rclk gpio.Bit
}

func NewShiftReg(din, clk, rclk iomux.Pin) *ShiftReg {
	din.Setup(iomux.D4mA)
	clk.Setup(iomux.D4mA)
	rclk.Setup(iomux.D4mA)
	gpio.UsePin(din)
	gpio.UsePin(clk)
	gpio.UsePin(rclk)
	sr := &ShiftReg{
		din:  gpio.BitForPin(din),
		clk:  gpio.BitForPin(clk),
		rclk: gpio.BitForPin(rclk),
	}
	sr.din.EnableOut()
	sr.clk.EnableOut()
	sr.rclk.EnableOut()
	return sr
}

func (sr *ShiftReg) WriteByte(b byte) {
	for i := 7; i >= 0; i-- {
		v := int(b) >> uint(i)
		sr.clk.Clear()
		sr.din.Store(v)
		sr.clk.Set()
	}
}

func (sr *ShiftReg) Latch() {
	sr.rclk.Set()
	sr.rclk.Clear()
}

// Segment bits
const (
	A = 1 << iota
	B
	C
	D
	E
	F
	G
	Q // Colon
)

// A TimeMux8Seg represents a time multiplexed 8-segment display. It assumes
// the display digits are driven by a 16-bit schift register (also cascade of
// two 8-bit registers). The 8 LSBits are connected to the anodes of all
// segments of the display and the 8 MSBits are connected to the common catodes
// of up to 8 supported digits.
type TimeMux8seg struct {
	digits [8]byte
	i, n   int8
	run    bool
}

var Display = NewTimeMux8seg(4, pins.GP11, pins.GP10, pins.GP9)

func NewTimeMux8seg(n int, din, clk, rclk iomux.Pin) *TimeMux8seg {
	if uint(n) > 8 {
		n = 8
	}
	d := new(TimeMux8seg)
	d.n = int8(n)
	go run(d, NewShiftReg(din, clk, rclk))
	return d
}

func (d *TimeMux8seg) SetSymbol(pos int, b byte) {
	if uint(pos) >= uint(d.n) {
		return
	}
	d.digits[pos] = b
}

func (d *TimeMux8seg) SetByte(pos int, b byte) {
	if uint(pos) >= uint(d.n) {
		return
	}
	d.digits[pos] = conv(b)
}

func (d *TimeMux8seg) Clear() {
	clear(d.digits[:])
}

func (d *TimeMux8seg) SetPos(pos int) {
	if uint(pos) >= uint(d.n) {
		pos = 0
	}
	d.i = int8(pos)
}

func (d *TimeMux8seg) WriteByte(b byte) (err error) {
	i := int(d.i)
	if b == '.' && i > 0 {
		if b = d.digits[i-1]; b&Q == 0 {
			d.digits[i-1] = b | Q
			return
		}
	}
	if b == '\n' || b == '\r' {
		d.i = 0
		return nil
	}
	if i < int(d.n) {
		d.digits[i] = conv(b)
		d.i = int8(i + 1)
	}
	return
}

func (d *TimeMux8seg) Write(p []byte) (int, error) {
	for _, b := range p {
		d.WriteByte(b)
	}
	return len(p), nil
}

func (d *TimeMux8seg) WriteString(s string) (int, error) {
	for i := 0; i < len(s); i++ {
		d.WriteByte(s[i])
	}
	return len(s), nil
}

func run(d *TimeMux8seg, sr *ShiftReg) {
	const fps = 60
	delay := time.Second / time.Duration(60*int(d.n))
	for i, n := uint(0), uint(d.n); ; {
		sr.WriteByte(^byte(1 << i))
		sr.WriteByte(d.digits[i])
		sr.Latch()
		if i++; i == n {
			i = 0
		}
		time.Sleep(delay)
	}
}

const digits = "" +
	"\x3f" + // A|B|C|D|E|F   -> 0
	"\x06" + // B|C           -> 1
	"\x5b" + // A|B|G|E|D     -> 2
	"\x4f" + // A|B|C|D|G     -> 3
	"\x66" + // F|G|B|C       -> 4
	"\x6d" + // A|F|G|C|D     -> 5
	"\x7d" + // A|F|E|D|C|G   -> 6
	"\x27" + // F|A|B|C       -> 7
	"\x7f" + // A|B|C|D|E|F|G -> 8
	"\x6f" //   A|B|C|D|F|G   -> 9

const letters = "" +
	"\x77" +
	"\x7c" +
	"\x58" +
	"\x5e" +
	"\x79" +
	"\x71" +
	"\x3d" +
	"\x74" +
	"\x30" +
	"\x1e" +
	"\x00" +
	"\x38" +
	"\x00" +
	"\x54" +
	"\x5c" +
	"\x73" +
	"\x67" +
	"\x50" +
	"\x6d" +
	"\x78" +
	"\x1c" +
	"\x00" +
	"\x00" +
	"\x00" +
	"\x6e"

var lette1rs = [...]byte{
	E | F | A | B | C | G, // A
	F | E | D | C | G,     // b
	G | E | D,             // c
	B | E | D | C | G,     // d
	A | F | E | D | G,     // E
	A | F | E | G,         // F
	A | F | E | D | C,     // G
	F | E | G | C,         // h
	F | E,                 // I
	E | D | C | B,         // J
	0,
	F | E | D, // L
	0,
	E | G | C,         // n
	G | E | D | C,     // o
	F | E | A | B | G, // P
	F | A | B | G | C, // q
	E | G,             // r
	A | F | G | C | D, // S
	F | E | D | G,     // t
	E | D | C,         // u
	0,
	0,
	0,
	F | G | B | C | D, // y
}

func conv(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		c = digits[c-'0']
	case 'a' <= c && c <= 'y':
		c = letters[c-'a']
	case 'A' <= c && c <= 'Y':
		c = letters[c-'A']
	case c == '-':
		c = G
	case c == '_':
		c = D
	case c == '=':
		c = D | G
	case c == '.':
		c = Q
	default:
		c = 0
	}
	return c
}
