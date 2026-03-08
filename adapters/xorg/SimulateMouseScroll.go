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

// Scroll Wheel
// 4 = dy down
// 5 = dy up
// 6 = dx right
// 7 = dx left

func to_abs_int(value int) int {
	if value < 0 {
		return -1 * value
	}
	return value
}

func SimulateMouseScroll(bridge *Bridge, x uint, y uint, dx int, dy int) error {

	if bridge.display != nil {

		target_x := C.int(x)
		target_y := C.int(y)

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

		C.XFlush(bridge.display)

		for d := 0; d < to_abs_int(dy); d++ {

			target_button := C.uint(4)

			if dy < 0 {
				target_button = C.uint(5)
			}

			C.XTestFakeButtonEvent(
				bridge.display,
				target_button,
				1,
				0,
			)

			C.XTestFakeButtonEvent(
				bridge.display,
				target_button,
				0,
				0,
			)

		}

		for d := 0; d < to_abs_int(dx); d++ {

			target_button := C.uint(6)

			if dx < 0 {
				target_button = C.uint(7)
			}

			C.XTestFakeButtonEvent(
				bridge.display,
				target_button,
				1,
				0,
			)

			C.XTestFakeButtonEvent(
				bridge.display,
				target_button,
				0,
				0,
			)

		}

		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
