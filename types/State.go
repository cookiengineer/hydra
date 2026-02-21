package types

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XInput2.h>
#include <X11/extensions/XTest.h>
#include <stdlib.h>
*/
import "C"

type State struct {
	Display        *C.Display
	XIOpcode       C.int
	XWindow        C.Window
	MouseEvents    chan MouseEvent
	KeyboardEvents chan KeyboardEvent
	VirtualScreen  *VirtualScreen
	running        bool
}

func (state *State) Destroy() {

	if state.Display != nil {

		C.XCloseDisplay(state.Display)
		state.Display = nil

	}

}

func (state *State) QueryPointer() (int, int, error) {

	if state.Display != nil {

		var returned_root     C.Window
		var returned_window   C.Window
		var returned_root_x   C.int
		var returned_root_y   C.int
		var returned_window_x C.int
		var returned_window_y C.int
		var returned_mask     C.uint

		result := C.XQueryPointer(
			state.Display,
			state.Window,
			&returned_root,
			&returned_window,
			&returned_root_x,
			&returned_root_y,
			&returned_window_x,
			&returned_window_y,
			&returned_mask,
		)

		if result == 0 {
			return 0, 0, errors.New("XQueryPointer failed")
		} else {
			return int(returned_root_x), int(returned_root_y), nil
		}

	} else {
		return 0, 0, errors.New("XDisplay is nil")
	}

}

// TODO: Continue Here
