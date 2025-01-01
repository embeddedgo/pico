#!/bin/sh

set -e

cd ../../../embeddedgo/pico/hal
hal=$(pwd)
cd ../p
rm -rf *

svdxgen github.com/embeddedgo/pico/p ../svd/*.svd

for p in iobank padsbank pllsys resets sio xosc; do
	cd $p
	xgen -g *.go
	GOOS=noos GOARCH=thumb $(emgo env GOROOT)/bin/go build -tags rp2350
	cd ..
done

exit

perlscript='
s/package irq/$&\n\nimport "embedded\/rtos"/;
s/ = \d/ rtos.IRQ$&/g;
'

cd $hal/irq
rm -f *
cp ../../p/irq/* .
perl -pi -e "$perlscript" *.go
