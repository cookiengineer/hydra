# Hydra LLM Bootstrap Document

## Architecture Overview

Hydra is a remote window manager that replaces Synergy/Barrier. It uses a one-binary-for-all design written in Go, targeting only X11 on Linux.

### Core Rules

1. **Controller machine** (the one running `hydra listen`) has the physical mouse and keyboard. All other machines are headless input-wise.
2. **Clients** receive input events from the controller via a long-lived HTTP streaming connection. Clients cannot modify server state.
3. **X11 only** and no Wayland support. Using XInput2 raw events and XTest for simulation.
4. **HTTP API** for debugging. SSH tunneling will be added later for security.
5. **/etc/hosts** and SSH keys are pre-configured on all machines.

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

### Key Bindings (processed on controller always)

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

#### Controller

| Binding            | Action                                          |
|:-------------------|:------------------------------------------------|
| `[Super]+[Escape]` | Reset to controller (deactivate remote machine) |

### State Machine

```
local-active: cursor on controller, events processed locally
- if mouse hits edge then activate remote machine, warp cursor
- if Super is held then check for key bindings (focus or tile)

remote-active: cursor on remote, events forwarded to remote socket
- if mouse returns to controller region then deactivate remote, warp cursor
- if Super+Escape then deactivate remote, warp cursor home
- local windows NEVER receive input while remote-active
```

### Protocol

JSON-over-HTTP, newline-delimited streaming.

**Endpoints:**

Server listens on `:3000` with JSON-over-HTTP newline-delimited streaming.

| Method | Path          | Purpose                                        |
|:-------|:--------------|:-----------------------------------------------|
| GET    | `/config`     | Returns full config (machines, virtual screen) |
| GET    | `/machines`   | Lists connected machines                       |
| POST   | `/connect`    | Connects client and enters event stream        |
| POST   | `/disconnect` | Disconnects client                             |

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

The server responds with a newline-delimited event stream:
Empty objects `{}` are 30-second keep-alive pings.

```
{"type":"init","virtual_screen":{"width":5760,"height":1080,"machines":{...}}}
{"type":"mouse","x":1920,"y":540,"dx":1,"dy":0,"button":0}
{"type":"keyboard","keycode":113}
{}
```


## Action Files

```bash
actions/
  Connect.go                   # client connects to controller, receives events
  Listen.go                    # controller event loop (signal handling, HTTP server, bridge init)
  WindowDrag.go                # experimental cross-machine window movement (process respawn)
  handleKeyboardEvent.go       # keyboard event routing (local passthrough vs remote forwarding)
  handleMouseEvent.go          # mouse event routing (boundary detection, remote forwarding)
  handleWindowManagerAction.go # key binding dispatch (focus navigation, tiling, reset)
```

## Helper Files

```bash
helpers/
  FindClosestWindowDown.go  # find the window below the focused window
  FindClosestWindowLeft.go  # find the window left of the focused window
  FindClosestWindowRight.go # find the window right of the focused window
  FindClosestWindowUp.go    # find the window above the focused window
```

## Adapter Files (xorg)

```bash
adapters/xorg/
  Bridge.go                # X11 display connection, event loop, pointer/modifier query
  Bridge_test.go           # SimulateMouseMove and QueryPointer tests
  FocusWindow.go           # XSetInputFocus + XRaiseWindow for window focusing
  HandleKeyboardEvent.go   # XInput raw key press/release → types.KeyboardEvent
  HandleMouseEvent.go      # XInput raw motion/button/scroll → types.MouseEvent
  MapWindow.go             # XMapWindow
  MoveResizeWindow.go      # XMoveResizeWindow, XMoveWindow, XResizeWindow
  QueryAllWindows.go       # XQueryTree → enumerate all visible top-level windows
  QueryFocusedWindow.go    # XGetInputFocus → get focused window with title/geometry
  RaiseWindow.go           # XRaiseWindow
  SimulateKeyboardEvent.go # dispatcher: routes to SimulateKeyPress/KeyRelease
  SimulateKeyPress.go      # XTestFakeKeyEvent (press)
  SimulateKeyRelease.go    # XTestFakeKeyEvent (release)
  SimulateMouseEvent.go    # dispatcher: routes to appropriate mouse simulation
  SimulateMouseMove.go     # XWarpPointer + XTestFakeMotionEvent
  SimulateMousePress.go    # XWarpPointer + XTestFakeButtonEvent (press)
  SimulateMouseRelease.go  # XWarpPointer + XTestFakeButtonEvent (release)
  SimulateMouseScroll.go   # fake scroll via button 4/5/6/7 press/release pairs
  TileWindow.go            # tile focused window to position on monitor
  UnmapWindow.go           # XUnmapWindow
```

## Type Definitions

```bash
types/
  Config.go                    # Config struct, machine management, virtual screen computation
  Config_test.go               # NewConfig, GetMachine, QueryMachine, Update/Remove/SetThis tests
  GlobalState.go               # thread-safe GlobalState (active machine, virtual screen)
  GlobalState_test.go          # SetActive, ResetActive, SetScreen tests
  IPv4.go                      # IPv4 type with validation, parsing, JSON serialization
  IPv4_test.go                 # IPv4 parse, IsIPv4, JSON round-trip tests
  IPv6.go                      # IPv6 type with validation, parsing, JSON serialization
  IPv6_test.go                 # IPv6 parse tests
  KeyBinding.go                # KeyBinding struct, modifier/key constants, action constants, Matches()
  KeyBinding_test.go           # match logic, default bindings, tile position tests
  KeyboardEvent.go             # KeyboardEvent struct (type + keycode)
  KeyboardEventType.go         # KeyPress / KeyRelease enum
  Machine.go                   # Machine struct (hostname, IP, position, screen, socket), Parse()
  Machine_test.go              # Parse valid/invalid positions, IPs, screens
  Monitor.go                   # Monitor struct (output, connected, resolution, geometry)
  MouseEvent.go                # MouseEvent struct (type, x, y, dx, dy, button)
  MouseEventButton.go          # MouseButtonLeft/Middle/Right enum
  MouseEventType.go            # MouseMove/MouseButtonPress/MouseButtonRelease/MouseScroll enum
  Screen.go                    # Screen struct (width, height, monitors, offsets)
  TilePosition.go              # TilePosition enum (8 positions), String(), FromString()
  VirtualScreen.go             # VirtualScreen struct (width, height, per-machine screens)
  Window.go                    # Window struct (id, title, x, y, width, height)
  computeVirtualScreen.go      # grid computation (left-to-right, top-to-bottom ordering)
  computeVirtualScreen_test.go # comprehensive grid layout tests
  formatIPv6.go                # IPv6 address normalization helper
  isHex.go                     # hex character check helper
```

## HTTP Handler Files

```
handlers/
  Handler_test.go # OnConnect validation tests, OnDisconnect tests
  OnConnect.go    # /connect endpoint: register machine, enter event stream
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
  ApplyKeyboardEvent.go      # keyboard event simulation on client (nil-safe bridge guard)
  ApplyKeyboardEvent_test.go # nil bridge noop and KeyPress/KeyRelease type tests
  ApplyMouseEvent.go         # mouse event simulation on client (global→local coordinate translation)
  ApplyMouseEvent_test.go    # coordinate translation with various VirtualScreen offsets
```

## CLI Entry Point

```
cmds/hydra/
  main.go # CLI: hydra listen <host>, hydra connect <position> <host>
```
