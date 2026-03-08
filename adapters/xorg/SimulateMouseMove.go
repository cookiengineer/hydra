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

func SimulateMouseMove(bridge *Bridge, x uint, y uint, dx int, dy int) error {

	if bridge.display != nil {

		origin_x := C.int(x)
		origin_y := C.int(x)
		target_x := C.int(int(x) + int(dx))
		target_y := C.int(int(y) + int(dy))

		if bridge.Screen != nil {
			offset_x := bridge.Screen.OffsetX
			offset_y := bridge.Screen.OffsetY
			origin_x = C.int(int(x) - int(offset_x))
			origin_y = C.int(int(y) - int(offset_y))
			target_x = C.int(int(x) + int(dx) - int(offset_x))
			target_y = C.int(int(y) + int(dy) - int(offset_y))
		}

		C.XWarpPointer(
			bridge.display,
			0,
			bridge.window,
			0, 0, 0, 0,
			origin_x,
			origin_y,
		)

		C.XTestFakeMotionEvent(
			bridge.display,
			-1,
			target_x,
			target_y,
			C.CurrentTime,
		)

		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}


