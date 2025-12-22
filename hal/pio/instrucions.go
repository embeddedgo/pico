// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pio

// JMP cond
const (
	Always int8 = 0
	Xzero  int8 = 1
	XnzDec int8 = 2
	Yzero  int8 = 3
	YnzDec int8 = 4
	XneqY  int8 = 5
	PIN    int8 = 6
	OSRne  int8 = 7
)

func JMP(newPC int, cond int8, delaySideSet uint32) uint32 {
	return 0 |
		uint32(newPC)&31 | uint32(cond)&7<<5 | delaySideSet&31<<8
}

// MOV src and dst
const (
	PINS uint8 = 0
	X    uint8 = 1
	Y    uint8 = 2
	EXEC uint8 = 4
	ISR  uint8 = 6
	OSR  uint8 = 7

	PINDIRS uint8 = 3 // destination only
	NULL    uint8 = 3 // source only

	PC     uint8 = 5 // destination only
	STATUS uint8 = 5 // source only
)

// MOV op
const (
	None       int8 = 0
	Invert     int8 = 1
	BitReverse int8 = 2
)

func MOV(dst uint8, op int8, src uint8, delaySideSet uint32) uint32 {
	return 0b101<<13 |
		uint32(src)&7 | uint32(op)&3<<3 | uint32(dst)&7<<5 | delaySideSet&31<<8
}

func btou32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func PUSH(ifEmpty, block bool, delaySideSet uint32) uint32 {
	return 0b100<<13 |
		btou32(block)<<5 | btou32(ifEmpty) | delaySideSet&31<<8
}

func PULL(ifEmpty, block bool, delaySideSet uint32) uint32 {
	return 0b100<<13 | 1<<7 |
		btou32(block)<<5 | btou32(ifEmpty) | delaySideSet&31<<8
}
