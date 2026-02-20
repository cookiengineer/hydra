package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XInput2.h>
#include <stdlib.h>
*/
import "C"

import "errors"
import "github.com/cookiengineer/hydra/types"

type State struct {
	XDisplay       *C.Display
	XIOpcode       C.int
	XWindow        C.Window
	MouseEvents    chan types.MouseEvent
	KeyboardEvents chan types.KeyboardEvent
	running        bool
}

func (state *State) Destroy() {

	if state.XDisplay != nil {
		C.XCloseDisplay(state.XDisplay)
		state.XDisplay = nil
	}

}

func (state *State) QueryPointer() (int, int, error) {

	var rootReturn C.Window
	var childReturn C.Window

	var rootX C.int
	var rootY C.int
	var winX C.int
	var winY C.int

	var maskReturn C.uint

	res := C.XQueryPointer(
		state.XDisplay,
		state.XWindow,
		&rootReturn,
		&childReturn,
		&rootX,
		&rootY,
		&winX,
		&winY,
		&maskReturn,
	)

	if res == 0 {
		return 0, 0, errors.New("XQueryPointer failed")
	}

	return int(rootX), int(rootY), nil

}

func (state *State) WarpPointer(x, y int) error {

	if state.XDisplay == nil {
		return errors.New("XDisplay is nil")
	}

	C.XWarpPointer(
		state.XDisplay,
		0,             // src_window (0 = none)
		state.XWindow, // dest_window (root window)
		0, 0, 0, 0,    // src rectangle ignored
		C.int(x),      // dest x
		C.int(y),      // dest y
	)

	C.XFlush(state.XDisplay)

	return nil

}
