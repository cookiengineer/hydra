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
import "github.com/cookiengineer/hydra/types"

func SimulateMousePress(bridge *Bridge, x uint, y uint, button types.MouseEventButton) error {

	if bridge.display != nil {

		target_x := C.int(x)
		target_y := C.int(y)
		target_button := C.uint(uint(button))

		if bridge.Screen != nil {
			offset_x := bridge.Screen.OffsetX
			offset_y := bridge.Screen.OffsetY
			target_x = C.int(int(x) - int(offset_x))
			target_y = C.int(int(y) - int(offset_y))
		}

		C.XWarpPointer(
			bridge.display,
			0,
			bridge.window,
			0, 0, 0, 0,
			target_x,
			target_y,
		)

		C.XTestFakeButtonEvent(
			bridge.display,
			target_button,
			1,
			0,
		)

		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
