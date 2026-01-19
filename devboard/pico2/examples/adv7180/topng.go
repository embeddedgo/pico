// Copyright 2025 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
)

func fatalErr(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}

func main() {
	f, err := os.Open("picocom.log")
	fatalErr(err)
	defer f.Close()

	img := image.NewYCbCr(image.Rect(0, 0, 720, 220), image.YCbCrSubsampleRatio422)

	scanner := bufio.NewScanner(f)
	start := false
	var y, c, n int
	for scanner.Scan() {
		line := scanner.Text()
		if !start {
			if strings.HasPrefix(line, "0: ") {
				start = true
				y = 0
				c = 0
				n++
				clear(img.Y)
				clear(img.Cb)
				clear(img.Cr)
			} else {
				continue
			}
		} else if strings.HasPrefix(line, "219: ") {
			start = false
		}
		line = line[strings.IndexByte(line, ':'):][10:]
		for i := 0; i < len(line); i += 8 {
			//fmt.Println(i, i+8, y, c)
			u64, _ := strconv.ParseUint(line[i:i+8], 16, 32)
			byry := uint32(u64)
			img.Cb[c] = byte(byry >> 24)
			img.Y[y] = byte(byry >> 16)
			img.Cr[c] = byte(byry >> 8)
			img.Y[y+1] = byte(byry)
			y += 2
			c += 1
		}
		if start == false {
			f, err := os.Create(fmt.Sprintf("%03d.png", n))
			fatalErr(err)
			fatalErr(png.Encode(f, img))
			fatalErr(f.Close())
		}
	}
}
