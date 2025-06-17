#!/bin/sh

name=$(basename $(pwd))

picotool load $name.elf
picotool reboot
