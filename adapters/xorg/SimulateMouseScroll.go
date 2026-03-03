package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XTest.h>
#include <stdlib.h>
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

func SimulateMouseScroll(bridge *Bridge, dx int, dy int) error {

	if bridge.display != nil {

		C.XFlush(bridge.display)

		for d := 0; d < to_abs_int(dy); d++ {

			button := 4

			if dy < 0 {
				button = 5
			}

			C.XTestFakeButtonEvent(bridge.display, C.uint(button), 1, 0)
			C.XTestFakeButtonEvent(bridge.display, C.uint(button), 0, 0)

		}

		for d := 0; d < to_abs_int(dx); d++ {

			button := 6

			if dx < 0 {
				button = 7
			}

			C.XTestFakeButtonEvent(bridge.display, C.uint(button), 1, 0)
			C.XTestFakeButtonEvent(bridge.display, C.uint(button), 0, 0)

		}

		C.XFlush(bridge.display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
