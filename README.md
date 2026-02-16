
# hydra

:construction: Experimental Prototype :construction:

Extreme multi head monitor management server and client that integrates well
across multiple SSH hosts and aims to replace barrier's messed up configurations
in favor of using `pactl`, `ssh`, xorg`, and `xrandr`.

## Features

These are planned features, and for now not ready for production yet.

- [ ] `hydra listen <host>` listens to mouse and keyboard input events on a local machine
- [ ] `hydra connect left-of <host>` connects to a remote listener machine
- [ ] `hydra connect right-of <host>` connects to a remote listener machine

- [ ] `hydra exec <host> <program> <args>` executes a program on a remote machine
- [ ] `hydra open <host> <url>` executes `xdg-open` on a remote machine
- [ ] `hydra copy <host> <file>` copies a file to a remote machine

## Opinions

- The listener machine has mouse, keyboard, and audio devices connected
- The listener machine uses `pactl` (pulseaudio) to configure audio devices
- All machines use `xrandr` to configure monitors

## Building

This project uses `CGo` to link against `X11` and `Xi`, which are the packages
`libx11` and `libxi` on most distributions. That sadly can't be avoided.

```bash
# Install dependencies
sudo pacman -S go libx11 libxi;

# Build program
bash build.sh;

# Show usage help
./hydra help;
```

## Usage

The assumed setup relies on local networking to be configured, so that hostnames
are locally reachable and configured in the `/etc/hosts` file.

**Example**:

Let's assume two machines with each having multiple monitors configured via `xrandr`.
In this example, the left machine has the keyboard and mouse connected to it, whereas
the right machine listens to keyboard and mouse events coming from the left machine.

```bash
# cat /etc/hosts;
192.168.0.12 machine-left
192.168.0.10 machine-right
```

```bash
# On left machine with mouse and keyboard
hydra listen left-machine;

# On right machine
hydra connect right-of left-machine;
```

The helper methods allow to remote-control a listening hydra instance, so that
key binding integration like volume up/volume down etc can be integrated properly,
even when you're using them on a remote machine (from the keyboard's point of view).

## License

AGPL-3.0
