# Hydra Setup Guide (hydraone + hydratwo)

This guide configures Hydra as a **user-specific systemd service** on two
machines:

```
                 virtual screen (5760 x 1080)
  +---------------------+---------------------+-----------------+
  |  hydraone (left)    |  hydraone (right)   |  hydratwo       |
  |  1920x1080          |  1920x1080          |  1920x1080      |
  |  HDMI-A-1 +0+0      |  HDMI-A-0 +1920+0   |  HDMI-A-1 +0+0  |
  +---------------------+---------------------+-----------------+
                           ^                    ^
                           |                    |
                192.168.0.11 hydraone   192.168.0.12 hydratwo
                keyboard + mouse         no local input
                (controller / server)    (client)
```

- **hydraone** (`192.168.0.11`) is the **controller/server**. It has the
  physical keyboard and mouse and runs `hydra listen`.
- **hydratwo** (`192.168.0.12`) is the **client**. It runs
  `hydra connect right-of` and receives input events from hydraone.

Hydra talks JSON-over-HTTP on TCP port `3333`. It does **not** require SSH
between the two machines; SSH is only used by the optional end-to-end test
toolchain (see [TESTING.md](./TESTING.md)).

> Both machines are X11 only. Wayland is not supported. `xrandr` must work on
> both machines.
>
> See [TODO.md](./TODO.md) for known issues and workarounds (for example the
> i3 key-grab overlap and the button-release handling).

## 1. Prerequisites (both machines)

Install the CGo build dependencies and `xrandr`.

**Arch Linux (hydraone):**

```bash
sudo pacman -S --needed go libx11 libxi libxtst xorg-xrandr;
```

**Debian / Ubuntu (hydratwo):**

```bash
sudo apt install golang-go libx11-dev libxi-dev libxtst-dev x11-xserver-utils;
```

Make sure the two machines can reach each other on the LAN and (optionally) add
each other to `/etc/hosts`. Adding the names is convenient, but the examples
below also work using raw IPs.

`/etc/hosts` on **both** machines:

```
192.168.0.11 hydraone
192.168.0.12 hydratwo
```

## 2. Build and install (both machines)

From a checkout of the repository, build natively for the host (the
prerequisites above install the required CGo libraries) and install it:

```bash
make native;
mkdir -p ~/.local/bin;
install -m 0755 build/hydra-native-linux-amd64 ~/.local/bin/hydra;
```

`make native` names the binary with Go's platform tags
(`hydra-native-<goos>-<goarch>`). To produce a portable binary for a specific
distribution instead (requires `podman`), run `make debian` or
`make archlinux` and copy `build/hydra-<distro>-linux-amd64` to the target.
Run `make help` for all targets.

Confirm `~/.local/bin` is on your `PATH`. The service units below use the
absolute path `%h/.local/bin/hydra`, so a missing `PATH` entry is fine.

### Remote install helper (recommended)

Instead of copying files and writing the unit by hand, run the installer from
the controller. It auto-detects the remote OS, distribution and architecture,
picks the matching `build/hydra-<distro>-<goos>-<goarch>` binary (building it
with `make <distro>` if needed), uploads it, and writes/enables the systemd user
service:

```bash
# client on hydratwo, run from the repository root on hydraone:
go run ./toolchain/install.go hydratwo;
```

When the target is the local machine (its hostname matches, or the argument is
one of this machine's IPs or `localhost`), the installer runs **locally without
SSH** and defaults the role to `server`. So the controller can be installed the
same way:

```bash
# controller on hydraone (this machine):
go run ./toolchain/install.go hydraone;
```

The client target defaults to `-controller <this host>` (`hydraone`) and
`-address <this host's LAN IP>`. Useful flags:

| Flag                 | Purpose                                              |
|:---------------------|:-----------------------------------------------------|
| `-role client\|server` | client (default remote) or controller (default local) |
| `-controller <name>` | controller name used in the handshake                |
| `-address <ip>`      | address the client connects to                       |
| `-position <pos>`    | client position (`right-of`, default)                |
| `-port <port>`       | `HYDRA_PORT` override                                |
| `-user <name>`       | SSH user (default: current user)                     |
| `-identity <path>`   | SSH private key (default: agent, then `~/.ssh/id_*`) |
| `-dry-run`           | detect and print the unit without changing anything  |
| `-verify-host-key`   | verify the host against `~/.ssh/known_hosts`         |

SSH authentication uses `ssh-agent` when available, otherwise the default
`~/.ssh/id_ed25519`, `id_ecdsa` and `id_rsa` keys. Each candidate key is tried
in its own connection so an unauthorized key does not abort the others.

## 3. Controller setup on hydraone

Create the systemd user unit:

```bash
mkdir -p ~/.config/systemd/user;
cat > ~/.config/systemd/user/hydra.service <<'EOF'
[Unit]
Description=Hydra controller (hydraone)
Documentation=https://github.com/cookiengineer/hydra
After=graphical-session.target
PartOf=graphical-session.target

[Service]
Type=simple
Environment=DISPLAY=:0
Environment=XAUTHORITY=%h/.Xauthority
ExecStart=%h/.local/bin/hydra listen hydraone
Restart=on-failure
RestartSec=2

[Install]
WantedBy=graphical-session.target
EOF
```

Enable and start it:

```bash
systemctl --user daemon-reload;
systemctl --user enable --now hydra.service;
systemctl --user status hydra.service;
```

Follow the logs:

```bash
journalctl --user -u hydra.service -f;
```

On startup Hydra writes its default configuration to
`~/.config/hydra/config.json` and begins listening on `:3333`. You do not need
to edit the file for the setup described here.

The default port is `3333`. Override it on **both** machines by setting
`HYDRA_PORT` in the unit (for example `Environment=HYDRA_PORT=4444`); the server
binds to it and the client connects to it.

### Open the firewall (only if one is active)

The client connects to the controller on TCP port `3333`:

```bash
# ufw
sudo ufw allow 3333/tcp;

# nftables (example, adapt to your ruleset)
sudo nft add rule inet filter input tcp dport 3333 accept;
```

If no firewall is active, no action is needed.

## 4. Client setup on hydratwo

The client must tell the controller its own name during the handshake. Because
it connects by IP address, set `HYDRA_CONTROLLER=hydraone` so the handshake
header matches the controller's configured name.

Create the systemd user unit:

```bash
mkdir -p ~/.config/systemd/user;
cat > ~/.config/systemd/user/hydra.service <<'EOF'
[Unit]
Description=Hydra client (hydratwo)
Documentation=https://github.com/cookiengineer/hydra
After=graphical-session.target
PartOf=graphical-session.target

[Service]
Type=simple
Environment=DISPLAY=:0
Environment=XAUTHORITY=%h/.Xauthority
Environment=HYDRA_CONTROLLER=hydraone
ExecStart=%h/.local/bin/hydra connect right-of 192.168.0.11
Restart=always
RestartSec=3

[Install]
WantedBy=graphical-session.target
EOF
```

Enable and start it:

```bash
systemctl --user daemon-reload;
systemctl --user enable --now hydra.service;
systemctl --user status hydra.service;
journalctl --user -u hydra.service -f;
```

`Restart=always` makes the client reconnect automatically if the controller
restarts or the network drops.

If you added `hydratwo` and `hydraone` to `/etc/hosts`, you can optionally
connect by name instead and drop the `HYDRA_CONTROLLER` line:

```ini
ExecStart=%h/.local/bin/hydra connect right-of hydraone
```

## 5. If `graphical-session.target` is not started by your session

Some i3/X11 sessions never activate `graphical-session.target`, which means a
unit installed with `WantedBy=graphical-session.target` will not start
automatically. Pick one of these fallbacks:

**Option A - start it from i3** (recommended for i3 setups):

Add to `~/.config/i3/config` on both machines:

```
exec_always --no-startup-id systemctl --user restart hydra.service
```

**Option B - bind to `default.target` instead:**

Change the `[Install]` section and reload:

```ini
[Install]
WantedBy=default.target
```

```bash
systemctl --user daemon-reload;
systemctl --user reenable hydra.service;
```

With Option B the service starts with the user manager (i.e. at login). If the
X server is not ready yet, `Restart=on-failure` / `Restart=always` and
`RestartSec` retry until it is.

`loginctl enable-linger $USER` is **not** required for a graphical service and
is not recommended here, because there is no X display before you log in.

## 6. Usage

Once both services are running, move the mouse to the **right edge of
hydraone's right monitor** to cross onto hydratwo. Move it back to hydratwo's
**left edge** to return to hydraone. The controller warps the local pointer and
restores the previously focused window on return.

Key bindings are processed on the controller and routed to whichever machine is
active (`Super` = Mod4):

| Binding                  | Action                                              |
|:-------------------------|:----------------------------------------------------|
| `Super`+`Left/Right/Up/Down` | Focus window in that direction (loops across machines) |
| `Super`+`Shift`+`Left/Right/Up/Down` | Tile focused window to that half          |
| `Super`+`Escape`         | Reset to the controller                             |
| `Super`+`~`/`0..9`/`-`/`+`/`Backspace` | Switch workspace (14 workspaces)      |
| `Super`+`Shift`+ same    | Move focused window to that workspace               |

Focus navigation is a loop across both machines: `Super`+`Right` from
hydraone enters hydratwo's left-most window, and continues off hydratwo's right
edge back to hydraone's left-most window. `Super`+`Left` is the reverse.

To change the client's position on the virtual screen, edit the `connect`
argument (`left-of`, `right-of`, `above`, `below`) and restart the client:

```bash
systemctl --user restart hydra.service;
```

### Coexisting with i3

Both machines here run i3 with the same config, so hydra's `Super`+arrow
bindings overlap with i3's own focus/tiling bindings.

- While a remote machine is active, hydra clears the local X input focus
  (`XSetInputFocus(None)`), so the focused local window does **not** receive
  typing. The last focused controller window is restored when control returns
  (edge crossing, `Super`+`Escape`, disconnect, or service shutdown).
- Because i3 uses key grabs for its own `Super`+key bindings, those bindings can
  still fire locally while a remote is active (only normal typing is
  suppressed). If that is undesirable, remove or remap the relevant
  `bindsym $mod+...` lines in the i3 config and let hydra own them.
- While the controller is active, hydra and i3 both react to `Super`+arrows.
  Hydra's geometric navigation and i3's tree navigation usually land on the
  same window.
- Hydra records the last focused controller window and restores it on return and
  on clean shutdown. If an unclean kill ever leaves focus odd, `i3-msg restart`
  on the affected machine re-synchronises i3.

## 7. Verification

On **hydraone**, query the controller API:

```bash
curl -s http://127.0.0.1:3333/machines | python3 -m json.tool;
curl -s http://127.0.0.1:3333/config  | python3 -m json.tool;
```

`/machines` should list both `hydraone` and `hydratwo`, and `/config` should
report a virtual screen that is the sum of both widths (here `5760x1080`).

Check the journal on either machine:

```bash
journalctl --user -u hydra.service -e;
```

Typical successful output:

```
Client connected: hydratwo (192.168.0.12)
```

## 8. Troubleshooting

| Symptom | Cause / fix |
|:--|:--|
| `Cannot open X display` | `DISPLAY`/`XAUTHORITY` wrong, or the service started before X. Check `echo $DISPLAY` in your session and confirm `~/.Xauthority` exists. |
| Connection fails with HTTP `412` | Handshake mismatch. The client must set `HYDRA_CONTROLLER=<controller-name>` (here `hydraone`) when connecting by IP. |
| `xrandr` error | `xrandr` is missing or the X server is not reachable. Install `xorg-xrandr` / `x11-xserver-utils`. |
| Client connects, then disconnects in a loop | Controller not reachable on port `3333`, or firewall blocks it. Try `curl http://192.168.0.11:3333/machines` from hydratwo. |
| `listen tcp :3333: bind: address already in use` | Another process already uses port `3333`. Stop it or change the port. |
| Workspaces/tiling stop responding | A remote client disconnected. Use `Super`+`Escape` to reset to the controller, or restart the service. |
| Window focus feels broken after killing hydra | Restart i3 with `i3-msg restart`, then restart `hydra.service`. A clean `systemctl --user stop` restores focus automatically. |
| A key or mouse button stays held on the remote | Stops when control returns to the controller. If it persists, restart the client service so it re-syncs. |

## 9. Security

Hydra's HTTP protocol is **unencrypted** and has no authentication beyond a
protocol/name handshake. Run it only on a trusted LAN. SSH tunnelling is on the
roadmap; until then, restrict port `3333` to your machines with a firewall.

## 10. Updating

```bash
git pull;
make native;
install -m 0755 build/hydra-native-linux-amd64 ~/.local/bin/hydra;
systemctl --user restart hydra.service;
```

Run this on each machine (restart the controller last so the client reconnects
to the new build).
