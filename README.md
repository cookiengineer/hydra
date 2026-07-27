
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
- [x] Virtual screen grid with tiling support
- [x] Keyboard-driven tiling (Super+Arrow keys)
- [x] i3-style workspaces (14 workspaces, Super+top-row keys)
- [x] X11 window management (move, resize, raise, focus, map/unmap)
- [x] HTTP API for programming and debugging (`/config`, `/connect`, `/disconnect`, `/machines`)
- [ ] SSH tunnel security (planned)
- [ ] Remote audio integration with pulseaudio (planned)
- [ ] Remote clipboard augmentation (e.g. `file://` links become `ssh://remote-host` links)
- [ ] Programmable Window Manager, so that AI assistants can interact with multi-head remote machines

## Opinions

- The controller machine has mouse, keyboard, and audio devices connected
- All machines use `xrandr` to configure monitors
- Machines are configured via `/etc/hosts` and have mutual SSH key access

## Building

This project uses CGo to link against `X11`, `Xi`, and `Xtst`.

### Arch Linux

```bash
sudo pacman -S go libx11 libxi;
bash build.sh;
```

### Debian / Ubuntu

```bash
sudo apt install golang-go libx11-dev libxi-dev libxtst-dev;
bash build.sh;
```

### Run Tests

```bash
go test ./handlers/... ./helpers/... ./parsers/... ./receivers/... ./types/...
```

## Usage

The assumed setup relies on local networking to be configured, so that hostnames
are locally reachable and configured in the `/etc/hosts` file.

**Example** (`/etc/hosts`):

```
192.168.0.12 controller
192.168.0.10 laptop
192.168.0.11 comachine
```

**Starting the controller** (machine with mouse and keyboard):

```bash
# On controller (has physical keyboard and mouse)
hydra listen controller;
```

**Connecting clients**:

```bash
# On laptop (appears left of controller on virtual screen)
hydra connect left-of controller;

# On comachine (appears right of controller on virtual screen)
hydra connect right-of controller;
```

Virtual screen layout becomes:

```
[laptop] [controller] [comachine]
 1280px    1920px       1280px
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
| `[Super]+[0]`..`[9]`  | Switch to 1..10   |
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

## Documentation

- Use the [BOOTSTRAP.md](./docs/BOOTSTRAP.md) to teach your LLM how to extend and use hydra as a window manager.
- Read the [TESTING.md](./docs/TESTING.md) for details on the in-container integration test workflow.

## License

AGPL-3.0
