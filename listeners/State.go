package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
#include <stdlib.h>
*/
import "C"

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
