# Hydra Xorg Integration Testing

## Overview

The `adapters/xorg` package uses CGo to call `libX11`, `libXi`, and `libXtst`
directly. Integration tests verify keyboard simulation, mouse simulation, and
window manipulation end-to-end through the actual CGo bindings.

Integration tests run inside an **isolated Xvfb container** so they never
interfere with the developer's desktop session. XTest key/mouse events,
pointer warping, and input focus changes are all contained.

## Quick Start

This builds a Podman container with Xvfb, compiles the integration-tagged
tests, and runs them on a virtual framebuffer (`:99`).

```bash
bash test.sh
```

## Manual Usage with Container

```bash
podman build -t hydra-xorg-test -f Containerfile .;
podman run --rm hydra-xorg-test;
```

## Manual Usage on Development Host

If you have Xvfb installed locally and want to run without a container:

```bash
# Install Xvfb first
sudo pacman -S xorg-server-xvfb    # Arch
# or
sudo apt-get install xvfb          # Debian/Ubuntu

# Start Xvfb
Xvfb :99 -screen 0 1920x1080x24 &

# Run tests
export DISPLAY=:99;
go test -tags=integration -v -count=1 ./adapters/xorg/...;
```

## Test Architecture

- Tests use `//go:build integration` tag, so they're excluded from `go test ./...` by default
- Tests use the `DISPLAY=":99"` inside the container without any window manager
- Tests use the `XCreateSimpleWindow` binding to create windows
- Tests share the same C types (`C.Display`, `C.Window`)

