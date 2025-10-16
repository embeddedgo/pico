#!/bin/sh

set -e

cd ../../../embeddedgo/pico/hal
hal=$(pwd)
cd ../p
rm -rf *

svdxgen github.com/embeddedgo/pico/p ../svd/*.svd

for p in clocks dma i2c iobank padsbank pio pll qmi resets sio spi ticks uart xosc; do
	cd $p
	xgen -g *.go
	GOTOOLCHAIN=go1.24.5-embedded GOOS=noos GOARCH=thumb go build -tags rp2350
	cd ..
done

perlscript='
s/package irq/$&\n\nimport "embedded\/rtos"/;
s/ = \d/ rtos.IRQ$&/g;
'

cd $hal/irq
rm -f *
cp ../../p/irq/* .
perl -pi -e "$perlscript" *.go
