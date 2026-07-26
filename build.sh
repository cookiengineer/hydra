#!/bin/bash

# Arch Linux:
#   sudo pacman -S go libx11 libxi
#
# Debian / Ubuntu:
#   sudo apt install golang-go libx11-dev libxi-dev libxtst-dev

go build -x -o ./hydra ./cmds/hydra/main.go;
