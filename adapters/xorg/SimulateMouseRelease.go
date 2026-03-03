package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
#include <X11/extensions/XTest.h>
*/
import "C"

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateMouseRelease(bridge *Bridge, button types.MouseEventButton) error {

	if bridge.display != nil {

		C.XTestFakeButtonEvent(bridge.display, C.uint(uint(button)), 0, 0)
		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
