package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
#include <X11/extensions/XTest.h>
*/
import "C"

import "errors"

func MoveResizeWindow(bridge *Bridge, window_id uint64, x int, y int, width int, height int) error {

	if bridge.display == nil {
		return errors.New("Display is nil")
	}

	C.XMoveResizeWindow(
		bridge.display,
		C.Window(window_id),
		C.int(x),
		C.int(y),
		C.uint(width),
		C.uint(height),
	)

	C.XFlush(bridge.display)

	return nil

}

func MoveWindow(bridge *Bridge, window_id uint64, x int, y int) error {

	if bridge.display == nil {
		return errors.New("Display is nil")
	}

	C.XMoveWindow(
		bridge.display,
		C.Window(window_id),
		C.int(x),
		C.int(y),
	)

	C.XFlush(bridge.display)

	return nil

}

func ResizeWindow(bridge *Bridge, window_id uint64, width int, height int) error {

	if bridge.display == nil {
		return errors.New("Display is nil")
	}

	C.XResizeWindow(
		bridge.display,
		C.Window(window_id),
		C.uint(width),
		C.uint(height),
	)

	C.XFlush(bridge.display)

	return nil

}

