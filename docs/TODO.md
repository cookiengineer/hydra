# Hydra TODO / Known Issues

Issues and rough edges collected while bringing hydra up on the `hydraone`
(controller) + `hydratwo` (client) pair.

Legend: `[ ]` open, `[~]` mitigated / worked around, `[x]` resolved.

## Known bugs

- [~] **Local i3 key grabs still fire while a remote is active.**
  The controller uses `XSetInputFocus(None)` to stop the focused window from
  receiving typing while a remote is active. However, i3 uses passive
  `XGrabKey` grabs for its own `Super`+key bindings, so those still fire
  locally. For example `Super`+`Return` can spawn a local terminal while the
  remote is active. Only normal typing is suppressed.
  An `XGrabKeyboard` attempt to fully block local input was reverted: it made
  keys land in *both* the local window and the remote, and/or broke XI2 raw key
  delivery. A reliable capture mechanism is still needed.

- [~] **`XSetInputFocus(None)` bypasses the window manager's focus model.**
  Hydra stores and restores the last focused controller window on return,
  `Super`+`Escape`, and clean shutdown, but an unclean kill (crash / `SIGKILL`)
  leaves X focus at `None` and i3 out of sync until `i3-msg restart` is run.
  A cooperative focus handling path (or always-restoring on every exit) would
  be more robust.

- [~] **Left mouse button release is swallowed by the implicit pointer grab.**
  On Xorg, `XI_RawButtonRelease` is not delivered to other clients after a
  button press (a lone release is delivered; a release following a press is
  not). This left the left button stuck on the remote, which made kitty hide
  its cursor. Workaround: a 100 ms `reconcileButtons` poll on the controller
  queries the physical button mask and synthesizes any missing release. The
  underlying X behaviour is not fixed; a release could still be missed if the
  button state changes faster than the poll.

- [~] **`Bridge.Init`'s event loop cannot be stopped cleanly.**
  `Destroy()` no longer closes the display while the `XNextEvent` goroutine is
  blocked (that caused a shutdown segfault), but this means the display is
  intentionally leaked on shutdown. A proper wakeup (self-pipe /
  `XAddConnectionWatch` / non-blocking loop with an exit flag) would allow a
  clean teardown.

- [ ] **i3 config overlap / double focus handling.**
  Both machines use the same i3 config, so hydra and i3 both bind
  `Super`+arrows. While the controller is active both act, which can
  double-move focus in some layouts. Either unbind `$mod+Left/Right/Up/Down`
  in i3 and let hydra own them, or make hydra defer to i3 while local-active.

- [ ] **`ApplyMouseEvent` has no nil guard for `virtualScreen`.**
  If a mouse event is processed before the `init` event populates the virtual
  screen, `virtualScreen.GetMachine(...)` panics. Add a nil check.

- [ ] **Modifier state is not reconstructed across machines.**
  Modifiers are forwarded as individual press/release keys; the modifier mask
  is sent but not applied. Modifiers held across activation, differing layouts,
  and `AltGr` / level-3 modifiers are not fully handled.

- [ ] **`localIP()` can pick the wrong interface.**
  `actions.Connect` uses the first non-loopback IPv4, which may be a docker /
  virtual interface. `toolchain/install.go` already uses a route-based method;
  unify them.

- [ ] **Window navigation operates on raw root children (i3 frames).**
  `QueryAllWindows` lists direct root children, so with i3 the "windows" are
  frames titled `[i3 con] container around 0x...` rather than real clients.
  Consider `_NET_CLIENT_LIST` / `_NET_ACTIVE_WINDOW` for reparenting WMs.

- [ ] **`computeVirtualScreen` ordering is simplistic.**
  Deterministic, but only handles the documented single-level left/right and
  above/below arrangements; mixed 2D layouts are not modelled.

- [ ] **Remote workspaces / `MoveWindowToWorkspace` are minimally tested.**
  Workspace state is per-machine and only lightly exercised; restored window
  geometry can drift.

## Testing gaps

- [ ] **E2E tests only cover connection/protocol.**
  `toolchain/e2e` verifies `/machines`, `/config`, and that the client receives
  `init`. Mouse edge crossing, keyboard isolation, button release, and the
  cross-machine focus loop were verified manually, not automatically.
- [ ] No automated tests for `reconcileButtons`, server-side
  `releaseTrackedInput`, or the client-side `ReleaseAll` safety net.
- [ ] `make integration` only runs the Debian Xvfb container; no Arch run.
- [ ] No CI configuration (unit / integration / e2e).

## UX / packaging

- [ ] Port is a hardcoded default (`3333`) with a `HYDRA_PORT` env override; no
  config-file support for port or listen address.
- [ ] Installer defaults to **disabled** host-key verification; `-verify-host-key`
  is opt-in.
- [ ] No `uninstall` / `update` installer subcommand.
- [ ] No log levels or structured logging; diagnostics use `fmt.Printf`.
- [ ] ARM support: the `Makefile` already suffixes `GOOS`/`GOARCH`, but there is
  no cross toolchain, so non-native builds are not yet possible.

## Security

- [ ] HTTP is unencrypted and unauthenticated beyond the
  `X-Protocol`/`X-Hydra-Controller` handshake; SSH tunnelling is planned.
- [ ] Any LAN host that knows the controller name can register a client.

## Forward-looking

- [ ] Remote audio integration (pulseaudio).
- [ ] Clipboard augmentation (e.g. `file://` links -> `ssh://remote-host`).
- [ ] Programmable window manager surface for AI assistants.

## Resolved

- [x] Broken `/connect` handshake (`X-Protocol` / `X-Hydra-Controller` /
  `RemoteAddr` IP parsing).
- [x] Event dispatch collision where `KeyPress == 0` / `MouseMove == 0` were
  dropped (introduced the `types.Event` envelope).
- [x] Raw-motion filter that discarded all axis-aligned mouse movement.
- [x] Relative mouse events treated as absolute (virtual cursor + absolute
  coordinates).
- [x] `KeyBinding.Matches` compared keysyms against keycodes, so shortcuts never
  fired (keysym -> keycode resolution).
- [x] Stuck modifier keys when control returned to the controller mid-hold
  (tracked-key release + client `ReleaseAll`).
- [x] Local pointer moved together with the remote cursor (server-side pointer
  pinning while remote-active).
- [x] Focus left at `None` after a clean service stop (restore on shutdown).
- [x] Shutdown segfault from closing the display under `XNextEvent` (guarded).
- [x] Broken cross-machine build/install (Makefile, per-distro Containerfiles,
  `toolchain/install.go`).
