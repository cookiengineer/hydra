
# hydra

Remote window manager that replaces Synergy/Barrier with a single Go binary.
Supports X11 on Linux (Arch, Debian, Ubuntu). Uses XInput2 raw events and
XTest for simulation.

## Features

- [x] `hydra listen <host>` listens for mouse/keyboard events (controller/server)
- [x] `hydra connect left-of <host>` connects to controller (position: left)
- [x] `hydra connect right-of <host>` connects to controller (position: right)
- [x] `hydra connect above <host>` connects to controller (position: above)
- [x] `hydra connect below <host>` connects to controller (position: below)
- [x] Virtual screen grid with absolute mouse handoff and edge crossing
- [x] Server-side virtual cursor (pointer keeps moving past the physical edge)
- [x] Keyboard-driven tiling (Super+Arrow keys)
- [x] Cross-machine `Super`+arrow focus loop with wraparound
- [x] i3-style workspaces (14 workspaces, Super+top-row keys)
- [x] X11 window management (move, resize, raise, focus, map/unmap)
- [x] Tracked input release (keys and mouse buttons released on return/disconnect)
- [x] HTTP API for programming and debugging (`/config`, `/connect`, `/disconnect`, `/machines`, `/event`)
- [x] Per-distribution builds (`debian`, `archlinux`) into `./build/`
- [x] Remote install helper (auto-detects remote OS/distro/arch, installs a systemd user service)
- [x] Two-machine end-to-end test toolchain (build, deploy over SSH, parse results)
- [ ] SSH tunnel security (planned)
- [ ] Remote audio integration with pulseaudio (planned)
- [ ] Remote clipboard augmentation (e.g. `file://` links become `ssh://remote-host` links)
- [ ] Programmable Window Manager, so that AI assistants can interact with multi-head remote machines

See [TODO](./docs/TODO.md) for the current list of known issues and
workarounds.

## Opinions

- The controller machine has mouse, keyboard, and audio devices connected
- All machines use `xrandr` to configure monitors
- All machines are configured via `/etc/hosts` and have mutual SSH key access

## Building

This project uses CGo to link against `X11`, `Xi`, and `Xtst`. Builds are
driven by the `Makefile`, which writes binaries into `./build/` using Go
platform tags in the file name (`hydra-<distro>-<goos>-<goarch>`).

Native build (uses your host's Go and dev libraries):

```bash
make native;   # -> build/hydra-native-linux-amd64
```

Distribution builds run inside `podman` containers, so the resulting binary is
portable to the matching distribution and to older glibc:

```bash
make debian;     # -> build/hydra-debian-linux-amd64
make archlinux;  # -> build/hydra-archlinux-linux-amd64
make all;        # both of the above
```

`GOOS`/`GOARCH` default to the host's Go values and can be overridden
(`make debian GOARCH=arm64`). Run `make help` for all targets.

### Remote Install

In order to install hydra and its systemd user service on another machine, run
the installer from the controller. It auto-detects the remote OS, distribution
and architecture and uploads the matching binary:

```bash
go run ./toolchain/install.go hydratwo;
```

If the target is this machine (its hostname or one of its IPs), the installer
runs locally without SSH and defaults to the `server` role:

```bash
go run ./toolchain/install.go hydraone;
```

### Run Tests

```bash
make unit;         # unit tests
make integration;  # Xvfb container integration tests
make e2e HOST=hydratwo;  # two-machine end-to-end tests
```

The end-to-end tests build a glibc-portable test binary in a `podman`
container, then copies it to the remote host over SSH, executes it there, and
parses the `stdout`/`stderr` results back on the controller.

See [docs/TESTING.md](./docs/TESTING.md) for details.

## Usage

The assumed setup relies on local networking to be configured, so that hostnames
are locally reachable and configured in the `/etc/hosts` file.

Example `/etc/hosts` file:

```
192.168.0.10 laptop
192.168.0.11 hydraone    # controller
192.168.0.12 hydratwo    # client / multi-head extension
192.168.0.13 hydrathree  # client / multi-head extension
# etc pp
```

Start the controller on the machine with mouse and keyboard:

```bash
# On controller (has physical keyboard and mouse)
hydra listen hydraone;
```

Then connect clients to the controller machine:

```bash
# On laptop (appears left of controller on virtual screen)
hydra connect left-of hydraone;

# On hydratwo (appears right of controller on virtual screen)
hydra connect right-of hydraone;
```

Virtual screen layout becomes:

```
[laptop] [hydraone] [hydratwo]
 1280px    1920px     1280px
```

## Window Manager Key Bindings

Key bindings are processed on the controller and work regardless of which machine is active.

### Focus Navigation

| Binding           | Action                    |
|:------------------|:--------------------------|
| `[Super]+[Left]`  | Focus window to the left  |
| `[Super]+[Right]` | Focus window to the right |
| `[Super]+[Up]`    | Focus window above        |
| `[Super]+[Down]`  | Focus window below        |

### Window Tiling

| Binding                   | Action                             |
|:--------------------------|:-----------------------------------|
| `[Super]+[Shift]+[Left]`  | Tile focused window to left half   |
| `[Super]+[Shift]+[Right]` | Tile focused window to right half  |
| `[Super]+[Shift]+[Up]`    | Tile focused window to top half    |
| `[Super]+[Shift]+[Down]`  | Tile focused window to bottom half |

### Controller

| Binding            | Action                                          |
|:-------------------|:------------------------------------------------|
| `[Super]+[Escape]` | Reset to controller (deactivate remote machine) |

### Workspace Switching (14 workspaces)

Workspaces are named FG, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, BG.

| Binding                | Action            |
|:-----------------------|:------------------|
| `[Super]+[~]`          | Switch to FG (0)  |
| `[Super]+[0]`..`[9]`   | Switch to 1..10   |
| `[Super]+[-]`          | Switch to 11      |
| `[Super]+[+]`          | Switch to 12      |
| `[Super]+[Backspace]`  | Switch to BG (13) |

### Workspace Moving

| Binding                      | Action                             |
|:-----------------------------|:-----------------------------------|
| `[Super]+[Shift]+[~]`        | Move focused window to FG (0)      |
| `[Super]+[Shift]+[0]`..`[9]` | Move focused window to 1..10       |
| `[Super]+[Shift]+[-]`        | Move focused window to 11          |
| `[Super]+[Shift]+[+]`        | Move focused window to 12          |
| `[Super]+[Shift]+[Backspace]`| Move focused window to BG (13)     |

## Status

Hydra runs as a pair of user systemd services on `hydraone` (controller) and
`hydratwo` (client). Mouse edge crossing, keyboard routing, tracked input
release, and the cross-machine `Super`+arrow focus loop work.

A number of rough edges around the concept remain. Notably i3 key-grab overlaps
with hydra's key grabs and the mouse button-release workaround.

See [TODO](./docs/TODO.md) for the current list.

## Documentation

- Read the [SETUP](./docs/SETUP.md) to install Hydra as a user systemd service on a controller and client machine.
- Read the [TODO](./docs/TODO.md) for known bugs, workarounds, and open work.
- Use the [BOOTSTRAP](./docs/BOOTSTRAP.md) to teach your LLM how to extend and use hydra as a window manager.
- Read the [TESTING](./docs/TESTING.md) for details on the container integration and two-machine end-to-end test workflows.

## License

AGPL-3.0
