package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
*/
import "C"

import "errors"

func SimulateMouseMove(bridge *Bridge, x int, y int) error {

	if bridge.display != nil {

		dest_x := C.int(x)
		dest_y := C.int(y)

		if bridge.Screen != nil {
			offset_x := bridge.Screen.OffsetX
			offset_y := bridge.Screen.OffsetY
			dest_x = C.int(x - offset_x)
			dest_y = C.int(y - offset_y)
		}

		C.XWarpPointer(
			bridge.display,
			0,
			bridge.window,
			0, 0, 0, 0,
			dest_x,
			dest_y,
		)

		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}


