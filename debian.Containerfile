# Debian build environment for hydra.
#
# Produces a Debian/glibc-portable binary that also runs on older Debian and
# Ubuntu targets (glibc <= the image's). The image is reused for the Xvfb
# integration tests, so it also carries xvfb.
#
# Build (via Makefile):
#   make debian
#
# Or manually:
#   podman build -t hydra-debian-build -f debian.Containerfile .

FROM docker.io/golang:bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    xvfb \
    libx11-dev \
    libxi-dev \
    libxtst-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
