package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XTest.h>
*/
import "C"

import "errors"

func SimulateKeyPress(bridge *Bridge, keycode uint32) error {

	if bridge.display != nil {

		C.XTestFakeKeyEvent(bridge.display, C.uint(keycode), 1, 0)
		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
