#!/bin/sh

name=$(basename $(pwd))

picotool load $name.uf2
