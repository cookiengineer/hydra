# Hydra LLM Bootstrap Document

## Architecture Overview

Hydra is a remote window manager that replaces Synergy/Barrier. It uses a one-binary-for-all design written in Go, targeting only X11 on Linux.

### Core Rules

1. **Controller machine** (the one running `hydra listen`) has the physical mouse and keyboard. All other machines are headless input-wise.
2. **Clients** receive input events from the controller via a long-lived HTTP streaming connection. Clients cannot modify server state.
3. **X11 only** and no Wayland support. Using XInput2 raw events and XTest for simulation.
4. **HTTP API** for debugging. The server listens on TCP `3333` (override with `HYDRA_PORT`). Client-to-controller requests send `X-Protocol: hydra` and `X-Hydra-Controller: <controller-name>`. SSH tunnelling will be added later for security.
5. **/etc/hosts** and SSH keys are pre-configured on all machines.

> See [TODO.md](./TODO.md) for known issues, workarounds, and open work. In
> particular, i3's own key grabs still fire locally while a remote is active,
> and the left mouse button release is reconstructed by a polling workaround.

### Config File

On first run, hydra auto-writes `~/.config/hydra/config.json` (XDG-compatible: `$XDG_CONFIG_HOME/hydra/config.json`). Contains:

- `workspaces`: 14 workspace definitions (Name + Index)
- `key_bindings`: customizable keyboard shortcuts including workspace bindings
- `controller`, `this`, `machines`, `screen`: runtime configuration

If the file is absent, defaults are written automatically.

### Virtual Screen Layout

```
[above-m2]  [above-m1]
[left-m2] [left-m1] [CONTROLLER] [right-m1] [right-m2]
[below-m1] [below-m2]
```

Each machine's screen is internally tiled. Supported tile positions:

- `left` which is the left half of active display
- `right` which is the right half of active display
- `top` which is the top half of active display
- `bottom` which is the bottom half of active display
- `top-left` which is the top-left quadrant
- `top-right` which is the top-right quadrant
- `bottom-left` which is the bottom-left quadrant
- `bottom-right` which is the bottom-right quadrant

### Key Bindings (processed on controller, routed to active machine)

Super key = Mod4 (left and right Super keys mapped to the same modifier).

#### Focus Navigation

| Binding           | Action                    |
|:------------------|:--------------------------|
| `[Super]+[Left]`  | Focus window to the left  |
| `[Super]+[Right]` | Focus window to the right |
| `[Super]+[Up]`    | Focus window above        |
| `[Super]+[Down]`  | Focus window below        |

#### Window Tiling

| Binding                   | Action                             |
|:--------------------------|:-----------------------------------|
| `[Super]+[Shift]+[Left]`  | Tile focused window to left half   |
| `[Super]+[Shift]+[Right]` | Tile focused window to right half  |
| `[Super]+[Shift]+[Up]`    | Tile focused window to top half    |
| `[Super]+[Shift]+[Down]`  | Tile focused window to bottom half |

#### Workspace Switching (14 workspaces)

| Binding                | Action              |
|:-----------------------|:--------------------|
| `[Super]+[~]`          | Switch to FG (0)    |
| `[Super]+[0]`..`[9]`  | Switch to 1..10     |
| `[Super]+[-]`          | Switch to 11        |
| `[Super]+[+]`          | Switch to 12        |
| `[Super]+[Backspace]`  | Switch to BG (13)   |

#### Workspace Moving

| Binding                      | Action                             |
|:-----------------------------|:-----------------------------------|
| `[Super]+[Shift]+[~]`        | Move focused window to FG (0)      |
| `[Super]+[Shift]+[0]`..`[9]` | Move focused window to 1..10       |
| `[Super]+[Shift]+[-]`        | Move focused window to 11          |
| `[Super]+[Shift]+[+]`        | Move focused window to 12          |
| `[Super]+[Shift]+[Backspace]`| Move focused window to BG (13)     |

Workspace names (in order): FG, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, BG.

#### Controller

| Binding            | Action                                          |
|:-------------------|:------------------------------------------------|
| `[Super]+[Escape]` | Reset to controller (deactivate remote machine) |

All key bindings (focus, tile, workspace, reset) route to whichever machine has input focus:

- if no remote active then applied locally on the controller
- if remote active then structured event sent to the remote via its event stream socket

### Workspaces

i3-style workspace management. Each workspace stores the geometry (X, Y, Width, Height) of
windows that belong to it.

- **Switch workspace**: save current layout, unmap all windows, restore target layout, map target windows
- **Move window**: remove focused window from current layout, add to target layout, unmap it
- Remote machines maintain their own workspace state, receiving `WorkspaceEvent` messages
- `GlobalState.ActiveWorkspace` is a string (workspace name, e.g. `"FG"`)
- `GlobalState.Workspaces` is `map[string]*Workspace` keyed by name

### State Machine

```
local-active: cursor on controller, events processed locally
- if mouse hits edge then remember the focused window, activate the remote
  machine, pin the local pointer to the boundary, send the initial absolute
  mouse position
- if Super is held then check for key bindings (focus, tile, workspace, reset)
- WM actions applied locally via X11 calls

remote-active: cursor on remote, events forwarded to remote socket
- the local pointer is pinned to the boundary each motion (XWarpPointer, which
  emits no raw events) so only the remote cursor moves
- raw device deltas accumulate into a virtual cursor; when it leaves the remote
  region on the controller's side, deactivate and warp the local pointer back
- pressed keys/buttons are tracked; they are released on deactivation, focus
  handoff, Super+Escape, disconnect, and clean shutdown
- a 100ms ticker reconciles the physical button mask and synthesizes any mouse
  button release missed because of the local implicit pointer grab
- if mouse returns to controller region then deactivate remote, warp cursor
- if Super+Escape then send ResetEvent to remote, deactivate remote, warp cursor
  home, restore last-focused controller window
- if Super+[workspace-key] then send WorkspaceEvent to remote
- if Super+[arrow] then send FocusEvent to remote; if the remote has no window
  in that direction it posts a handoff to /event and focus wraps to the
  controller's extreme window
- if Super+Shift+[arrow] then send TileEvent to remote
- on activation the controller calls XSetInputFocus(None) so normal typing does
  not reach local windows; the last focused window is restored on deactivation
  and on clean shutdown. Note: i3's own key grabs still fire locally while a
  remote is active (see TODO.md)
```

### Protocol

JSON-over-HTTP, newline-delimited streaming.

**Endpoints:**

Server listens on `:3333` with JSON-over-HTTP newline-delimited streaming.

| Method | Path          | Purpose                                        |
|:-------|:--------------|:-----------------------------------------------|
| GET    | `/config`     | Returns full config (machines, virtual screen) |
| GET    | `/machines`   | Lists connected machines                       |
| POST   | `/connect`    | Connects client and enters event stream        |
| POST   | `/disconnect` | Disconnects client                             |
| POST   | `/event`      | Client-to-controller events (e.g. focus handoff) |

All client-to-controller requests MUST send `X-Protocol: hydra` and
`X-Hydra-Controller: <controller-hostname>`. The controller compares the latter
(and the `Host` header as fallback, with the port stripped) against its own
configured controller name.

#### 1. Client Connect request format (client to server)

Clients send a POST to `/connect` with their machine info as JSON:

```
POST /connect
Content-Type: application/json
X-Protocol: hydra

{
  "hostname": "machine-bar",
  "ip": "192.168.0.10",
  "position": "left-of",
  "screen": {
    "width": 1920,
    "height": 1080,
    "monitors": [{"output": "HDMI-0", "width": 1920, "height": 1080}]
  }
}
```

#### 2. Event stream format (server to client)

The server responds with a newline-delimited event stream. Every line is a
JSON envelope with a string `type` and a `data` payload object. Empty objects
`{}` are 30-second keep-alive pings. The envelope exists because the numeric
`MouseEvent`/`KeyboardEvent` `type` fields use `0` as a valid value
(`MouseMove`, `KeyPress`), which previously made it impossible to tell an
event from an unset field.

```
{"type":"init","data":{"virtual_screen":{"width":5760,"height":1080,"machines":{...}},"workspaces":[...],"active_workspace":"FG"}}
{"type":"mouse","data":{"type":0,"x":1920,"y":540,"dx":0,"dy":0,"button":0}}
{"type":"keyboard","data":{"type":0,"keycode":114,"keysym":65363,"modifiers":64}}
{"type":"workspace","data":{"type":"workspace","name":"1","index":1}}
{"type":"focus","data":{"type":"focus","direction":"left"}}
{"type":"tile","data":{"type":"tile","position":"right"}}
{"type":"reset","data":{"type":"reset"}}
{}
```

Mouse events carry absolute virtual-screen coordinates. The controller
accumulates raw device deltas into a virtual cursor, so the pointer can keep
moving after the physical pointer reaches a screen edge. The client subtracts
its own screen offset before warping.

Keyboard events carry the keysym (not just the local keycode) plus the modifier
mask, so a client with a different keyboard layout maps the key correctly.

#### 3. Client-to-controller events (`POST /event`)

When the client has no window in the requested direction it hands focus back:

```
{"type":"handoff","machine":"hydratwo","direction":"right"}
```

The controller only accepts the handoff if `machine` is its currently active
machine. It then deactivates the remote and focuses its own extreme window
(`right` wraps to the left-most window, `left` to the right-most window).

#### 4. Event payloads

| Event Type    | Data schema                                             | Purpose                                  |
|:--------------|:--------------------------------------------------------|:-----------------------------------------|
| `init`        | `{"virtual_screen":...,"workspaces":[...],"active_workspace":"FG"}` | Client startup; carries workspace definitions |
| `workspace`   | `{"name":"1","index":1}`                                | Switch workspace on active machine       |
| `focus`       | `{"direction":"left"}` or `{"edge":"first"\|"last"}`    | Focus window in direction (or extreme) on remote |
| `tile`        | `{"position":"top"}`                                    | Tile focused window on remote            |
| `reset`       | `{}`                                                    | Unfocus remote's window; focus returns to controller |

`FocusEvent.Edge` selects the left-most (`first`) or right-most (`last`) window
and is used when entering a remote machine at a boundary.

The `init` event replaces the previous ad-hoc map payload. It includes the full workspace
definitions and the controller's current active workspace index so the client can initialize
its own `GlobalState`.


## Action Files

```bash
actions/
  Connect.go                   # client connects to controller, receives/dispatches events
  Listen.go                    # controller event loop (signal handling, HTTP server, bridge init,
                               # 100ms button reconcile ticker, clean shutdown focus restore)
  handleKeyboardEvent.go       # keyboard routing; tracks pressed keys, suppresses releases of
                               # key presses consumed by WM shortcuts
  handleMouseEvent.go          # mouse routing; virtual cursor accumulation, boundary detection,
                               # button tracking, server-side pointer pinning
  handleWindowManagerAction.go # key binding dispatch (focus, tile, workspace switch/move, reset)
                               # resolves keysyms to keycodes and routes structured events to remote
  remote.go                    # activation/deactivation, pointer pinning, tracked-input release,
                               # reliable keyboard/mouse sends, /event handoff handler, reconcileButtons()
  display.go                   # defaultDisplay() (DISPLAY env with :0 fallback)
```

## Helper Files

```bash
helpers/
  FindClosestWindowDown.go  # find the window below the focused window
  FindClosestWindowLeft.go  # find the window left of the focused window
  FindClosestWindowRight.go # find the window right of the focused window
  FindClosestWindowUp.go    # find the window above the focused window
  FindWindowEdge.go         # leftmost/rightmost/topmost/bottommost window (boundary handoff/wrap)
  SwitchWorkspace.go        # save current layout, unmap all windows, restore target layout
  MoveWindowToWorkspace.go  # move focused window to target workspace, unmap it
```

## Adapter Files (xorg)

```bash
adapters/xorg/
  Bridge.go                # X11 display connection (XInitThreads), event loop, pointer/modifier query
                           # atomic running flag; Destroy() leaks the display while the loop runs
  Bridge_test.go           # SimulateMouseMove and QueryPointer tests
  FocusWindow.go           # XSetInputFocus + XRaiseWindow for window focusing
  HandleKeyboardEvent.go   # XInput raw key press/release → types.KeyboardEvent
  HandleMouseEvent.go      # XInput raw motion/button/scroll → types.MouseEvent
                           # (ignores only zero-delta motion)
  Keysym.go                # XKeysymToKeycode / XkbKeycodeToKeysym translation
  MapWindow.go             # XMapWindow
  MoveResizeWindow.go      # XMoveResizeWindow, XMoveWindow, XResizeWindow
  QueryAllWindows.go       # XQueryTree → enumerate all visible top-level windows
  QueryFocusedWindow.go    # XGetInputFocus → get focused window with title/geometry
  RaiseWindow.go           # XRaiseWindow
  ReleaseModifiers.go      # XTestFakeKeyEvent release for common modifier keysyms
  SimulateKeyboardEvent.go # dispatcher: routes to SimulateKeyPress/KeyRelease
  SimulateKeyPress.go      # XTestFakeKeyEvent (press)
  SimulateKeyRelease.go    # XTestFakeKeyEvent (release)
  SimulateMouseEvent.go    # dispatcher: routes to appropriate mouse simulation
  SimulateMouseMove.go     # XWarpPointer + XTestFakeMotionEvent
  SimulateMousePress.go    # XWarpPointer + XTestFakeButtonEvent (press)
  SimulateMouseRelease.go  # XWarpPointer + XTestFakeButtonEvent (release)
  SimulateMouseScroll.go   # fake scroll via button 4/5/6/7 press/release pairs
  TileWindow.go            # tile focused window to position on monitor
  UnfocusWindow.go         # XSetInputFocus(None) — clears local focus while a remote is active
  UnmapWindow.go           # XUnmapWindow
```

## Type Definitions

```bash
types/
  Config.go                    # Config struct, machine management, virtual screen, workspace definitions, key bindings
  Config_test.go               # NewConfig, GetMachine, QueryMachine, Update/Remove/SetThis, workspace tests
  Event.go                     # Event envelope {type, data} + NewEvent()/Unmarshal()
  Event_test.go                # envelope round-trip and zero-value type tests
  FocusEvent.go                # FocusEvent (type, direction, edge) for remote focus routing
  FocusEvent_test.go           # JSON marshal/unmarshal round-trip tests
  GlobalState.go               # thread-safe GlobalState (active machine, virtual screen, workspaces,
                               # last focused window, virtual cursor, pressed keys/buttons)
  GlobalState_test.go          # SetActive, ResetActive, SetScreen, workspace state, layout storage tests
  GlobalStateKeys_test.go      # pressed key/button tracking tests
  HandoffEvent.go              # HandoffEvent {type, machine, direction} for POST /event
  IPv4.go                      # IPv4 type with validation, parsing, JSON serialization
  IPv4_test.go                 # IPv4 parse, IsIPv4, JSON round-trip tests
  IPv6.go                      # IPv6 type with validation, parsing, JSON serialization
  IPv6_test.go                 # IPv6 parse tests
  InitEvent.go                 # InitEvent struct (replaces ad-hoc map; carries virtual_screen, workspaces, active_workspace)
  InitEvent_test.go            # JSON marshal/unmarshal round-trip tests
  KeyBinding.go                # KeyBinding struct (modifiers, keycode, action, data), modifier/key/action constants, Matches(), GetDefaultKeyBindings()
  KeyBinding_test.go           # match logic, default bindings, workspace binding Data field tests
  KeyboardEvent.go             # KeyboardEvent struct (type, keycode, keysym, modifiers)
  KeyboardEventType.go         # KeyPress / KeyRelease enum
  Machine.go                   # Machine struct (hostname, IP, position, screen, socket), Parse()
  Machine_test.go              # Parse valid/invalid positions, IPs, screens
  Monitor.go                   # Monitor struct (output, connected, resolution, geometry)
  MouseEvent.go                # MouseEvent struct (type, x, y, dx, dy, button)
  MouseEventButton.go          # MouseButtonLeft/Middle/Right enum
  MouseEventType.go            # MouseMove/MouseButtonPress/MouseButtonRelease/MouseScroll enum
  ResetEvent.go                # ResetEvent struct (type) for remote reset routing
  ResetEvent_test.go           # JSON marshal/unmarshal round-trip tests
  Screen.go                    # Screen struct (width, height, monitors, offsets)
  TileEvent.go                 # TileEvent struct (type + position) for remote tile routing
  TileEvent_test.go            # JSON marshal/unmarshal round-trip tests
  TilePosition.go              # TilePosition enum (8 positions), String(), FromString()
  VirtualScreen.go             # VirtualScreen struct (width, height, per-machine screens)
  Window.go                    # Window struct (id, title, x, y, width, height)
  Workspace.go                 # Workspace struct (name, index, windows)
  Workspace_test.go            # Workspace field tests
  WorkspaceEvent.go            # WorkspaceEvent struct (type, name, index) for remote workspace routing
  WorkspaceEvent_test.go       # JSON marshal/unmarshal round-trip tests
  computeVirtualScreen.go      # grid computation (left-to-right, top-to-bottom ordering)
  computeVirtualScreen_test.go # comprehensive grid layout tests
  formatIPv6.go                # IPv6 address normalization helper
  isHex.go                     # hex character check helper
  Port.go                      # DefaultPort (3333) + Port() with HYDRA_PORT override
```

## HTTP Handler Files

```
handlers/
  Handler_test.go # OnConnect validation tests, OnDisconnect tests
  OnConnect.go    # /connect endpoint: register machine, send InitEvent, enter event stream
  OnDisconnect.go # /disconnect endpoint: remove machine, recompute virtual screen
```

## Parser Files

```
parsers/
  Xrandr.go      # parse xrandr --query output into types.Screen
  Xrandr_test.go # test parsing with mock xrandr output strings
```

## Receiver Files

```
receivers/
  ApplyFocusEvent.go          # client-side focus handler (closest in direction or edge); returns found/moved
  ApplyFocusEvent_test.go     # nil bridge guard, all directions
  ApplyKeyboardEvent.go       # keyboard simulation + pressed-keycode tracking (nil-safe bridge guard)
  ApplyKeyboardEvent_test.go  # nil bridge noop and KeyPress/KeyRelease type tests
  ApplyMouseEvent.go          # mouse simulation (global→local coordinate translation) + button tracking
  ApplyMouseEvent_test.go     # coordinate translation with various VirtualScreen offsets
  ApplyResetEvent.go          # client-side reset handler (unfocus + release modifiers)
  ApplyResetEvent_test.go     # nil bridge guard
  ApplyTileEvent.go           # client-side tile event handler (queries focused window, tiles on local monitor)
  ApplyTileEvent_test.go      # nil bridge guard, nil virtual screen guard, all positions
  ApplyWorkspaceEvent.go      # client-side workspace switch handler (calls helpers.SwitchWorkspace)
  ApplyWorkspaceEvent_test.go # nil bridge guard, nil state guard, all workspace names
  ReleaseAll.go               # releases all tracked keys/buttons + modifiers when the stream ends
  ReleaseAll_test.go          # tracking + reset tests
```

## CLI Entry Point

```
cmds/hydra/
  main.go # CLI: hydra listen <host>, hydra connect <position> <host>
```

## Toolchain (remote e2e tests)

```
toolchain/
  test.go        # main(): go run test.go <host> <server|client> [flags]
  install.go     # main(): go run install.go <host>  (//go:build install)
  lib/
    exec.go      # Executor interface + OSExecutor (command execution)
    plan.go      # Plan: pure command builders (podman build, scp, ssh run)
    parse.go     # parses --- PASS/FAIL and E2E_RESULT lines, Summary.OK()
    root.go      # locate repository root via go.mod
  e2e/
    e2e_test.go  # //go:build e2e; role-aware TestE2E (server/client)
```

`install.go` uses `golang.org/x/crypto/ssh` to detect the remote OS,
distribution and architecture, upload the matching `build/hydra-<distro>-...`
binary, and write/enable the systemd user service. When the target host is the
local machine (hostname, one of its IPs, or `localhost`) it runs locally
without SSH and defaults to the `server` role. It is tagged `//go:build
install` so `go build ./...` skips it while `go run install.go` still works; the
tag also keeps `golang.org/x/crypto` in `go.mod`. Run it from the controller:
`go run ./toolchain/install.go hydratwo` (client) or `... hydraone` (server).

The e2e builder reuses the shared `debian.Containerfile` (golang:bookworm +
libx11-dev libxi-dev libxtst-dev xvfb). The test binary is built in the
container for glibc portability, copied to the remote with `scp`, and executed
over SSH with `DISPLAY`/`XAUTHORITY`. Results are parsed from the combined
stdout/stderr. See `docs/TESTING.md`.

## Build System

```
Makefile                   # centralized build/test targets
debian.Containerfile       # Debian build environment (glibc-portable)
archlinux.Containerfile    # Arch Linux build environment
```

`make debian` / `make archlinux` / `make native` write binaries to `./build/`
named `hydra-<distro>-<goos>-<goarch>` (Go platform tags). `make unit`,
`make integration`, and `make e2e HOST=...` run the test suites.

`GlobalState` also tracks the virtual cursor (`SetCursor`/`GetCursor`) used for
absolute mouse positioning on the active machine, plus the keys and mouse
buttons currently forwarded to the active remote (`PressedKeys`,
`PressedButtons` and their accessors) so they can be released on deactivation,
handoff, disconnect, and shutdown.

## Known Issues

See [TODO.md](./TODO.md) for the current list of bugs and workarounds.

