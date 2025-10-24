#!/bin/sh

set -e

gdb_cmd=gdb
if command -v arm-none-eabi-gdb; then
	gdb_cmd=arm-none-eabi-gdb
elif command -v gdb-multiarch; then
	gdb_cmd=gdb-multiarch
fi

reset='monitor connect_srst enable'
if [ $# -eq 1 -a "$1" = 'noreset' ]; then
        reset=''
fi

remote='blackmagic.local:2345'
#remote='/dev/ttyACM0'

brkpnt=6
wchpnt=4

$gdb_cmd --tui \
	-ex "target extended-remote $remote" \
	-ex "$reset" \
	-ex 'monitor swdp_scan' \
	-ex 'attach 1' \
	-ex 'set mem inaccessible-by-default off' \
	-ex "set remote hardware-breakpoint-limit $brkpnt" \
	-ex "set remote hardware-watchpoint-limit $wchpnt" \
	-ex 'set history save on' \
	-ex 'set history filename ~/.gdb-history-embeddedgo' \
	-ex 'set history size 1000' \
	"$(basename $(pwd)).elf"
