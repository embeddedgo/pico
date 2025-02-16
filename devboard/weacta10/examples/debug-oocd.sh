#!/bin/sh

INTERFACE=cmsis-dap
SPEED=5000
TARGET=rp2350

. $(emgo env GOROOT)/../scripts/debug-oocd.sh
