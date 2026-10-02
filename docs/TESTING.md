# Hydra Xorg Integration Testing

## Overview

The `adapters/xorg` package uses CGo to call `libX11`, `libXi`, and `libXtst`
directly. Integration tests verify keyboard simulation, mouse simulation, and
window manipulation end-to-end through the actual CGo bindings.

Integration tests run inside an **isolated Xvfb container** so they never
interfere with the developer's desktop session. XTest key/mouse events,
pointer warping, and input focus changes are all contained.

## Quick Start

This builds the Debian Podman container (with Xvfb), compiles the
integration-tagged tests, and runs them on a virtual framebuffer (`:99`).

```bash
make integration;
# or: bash test.sh integration;
```

## Manual Usage with Container

```bash
podman build -t hydra-debian-build -f debian.Containerfile .;
podman run --rm -v "$PWD":/src:ro,Z -w /src -e DISPLAY=:99 hydra-debian-build \
    sh -c 'Xvfb :99 -screen 0 1920x1080x24 & sleep 1 && go test -tags=integration -v -count=1 ./adapters/xorg/...';
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

## End-to-End Testing (two machines)

The e2e workflow verifies a real controller/client pair over SSH. It lives in
`toolchain/`:

```
toolchain/test.go          # main() runner / orchestrator (go run test.go ...)
toolchain/lib/             # build, deploy, exec, and result-parsing helpers
toolchain/e2e/             # role-aware //go:build e2e tests
debian.Containerfile       # shared glibc-portable builder image
```

The e2e builder reuses `debian.Containerfile` (the same image family as
`make debian`).

### Why a container

The target host may not have the CGo build dependencies (`libxtst-dev`, and
`XTest.h`), and its glibc may be older than the developer host. The test binary
is therefore built inside a Debian-based `golang:bookworm` container and copied
to the target, where it links against the runtime `libX*.so.6` sonames.

### Usage

```bash
# From the repository root:
go run ./toolchain/test.go hydratwo server;   # build + deploy + orchestrate both ends
go run ./toolchain/test.go hydratwo client;   # run only the client role remotely
```

Common flags:

| Flag            | Default                        | Purpose                                   |
|:----------------|:-------------------------------|:------------------------------------------|
| `-controller`   | this hostname                  | controller name used for the handshake    |
| `-address`      | detected LAN IPv4              | address the client uses to reach hydraone |
| `-position`     | `right-of`                     | client position on the virtual screen     |
| `-test`         | `TestE2E`                      | `go test -run` filter                     |
| `-image`        | `hydra-e2e-builder`            | builder image name                        |
| `-display`      | `:0`                           | client `DISPLAY`                          |
| `-xauthority`   | `$HOME/.Xauthority`            | client `XAUTHORITY`                       |

Both hosts must run Xorg (`xrandr` must work) and the client must be reachable
via passwordless SSH. The tests temporarily bind `:3333` on the controller.

The workflow is scheduled through `test.sh e2e <host>`.

### Live verification notes

- The controller drives the remote cursor with absolute virtual-screen
  coordinates; raw device motion is accumulated server-side so the pointer can
  keep moving past the physical screen edge. The local pointer is pinned to the
  boundary with `XWarpPointer`, which emits no raw events.
- On activation the controller calls `XSetInputFocus(None)` so normal typing
  does not reach local windows; the previously focused window is restored on
  return, `Super`+`Escape`, disconnect, and clean shutdown.
- Pressed keys and mouse buttons are tracked and released on deactivation,
  focus handoff, disconnect, and shutdown. A 100 ms ticker also reconciles the
  physical button mask, because Xorg does not deliver `XI_RawButtonRelease`
  after a button press (see [TODO.md](./TODO.md)).
- `Super`+arrow focus navigation crosses machine boundaries and wraps: right
  from the controller enters the client's leftmost window, and right from the
  client's rightmost window wraps back to the controller's leftmost window.

### Current coverage

The e2e suite currently verifies **connection and protocol only**: the client
registers with `/connect`, receives `init`, and appears in `/machines` and
`/config`. Mouse edge crossing, keyboard isolation, tracked-input release, and
the cross-machine focus loop are verified manually. Expanding the e2e suite to
cover these is tracked in [TODO.md](./TODO.md).


