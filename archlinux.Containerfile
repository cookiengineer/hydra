# Arch Linux build environment for hydra.
#
# Produces a binary linked against Arch's rolling libraries. Prefer this build
# for Arch targets, and the debian.Containerfile build for Debian/Ubuntu.
#
# Build (via Makefile):
#   make archlinux
#
# Or manually:
#   podman build -t hydra-archlinux-build -f archlinux.Containerfile .

FROM docker.io/archlinux:latest

RUN pacman -Syu --noconfirm --needed \
    go \
    gcc \
    libx11 \
    libxi \
    libxtst \
    xorg-server-xvfb \
    && pacman -Scc --noconfirm

WORKDIR /src
