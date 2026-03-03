package listeners

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

import "errors"
import "github.com/cookiengineer/hydra/types"

type State struct {
	XDisplay       *C.Display
	XIOpcode       C.int
	XWindow        C.Window
	MouseEvents    chan types.MouseEvent
	KeyboardEvents chan types.KeyboardEvent
	VirtualScreen  *types.VirtualScreen
	running        bool
}


// -------------------- MOUSE SIMULATION --------------------

// dx, dy: scroll direction (1/-1)
func (state *State) SimulateMouseScroll(dx, dy int) {
	if state.XDisplay == nil {
		return
	}

	// Optional: clamp scrolling to local monitor
	var minX, minY, maxX, maxY int
	if state.VirtualScreen != nil && state.VirtualScreen.Active != nil {
		m := state.VirtualScreen.Active
		minX, minY = m.OffsetX, m.OffsetY
		maxX, maxY = m.OffsetX+m.Screen.Width-1, m.OffsetY+m.Screen.Height-1
		cx, cy, _ := state.QueryPointer()
		if cx < minX || cx > maxX || cy < minY || cy > maxY {
			return
		}
	}

	// proceed with scroll
	for i := 0; i < abs(dx); i++ {
		btn := 6
		if dx > 0 {
			btn = 7
		}
		C.XTestFakeButtonEvent(state.XDisplay, C.uint(btn), 1, 0)
		C.XTestFakeButtonEvent(state.XDisplay, C.uint(btn), 0, 0)
	}
	for i := 0; i < abs(dy); i++ {
		btn := 4
		if dy < 0 {
			btn = 5
		}
		C.XTestFakeButtonEvent(state.XDisplay, C.uint(btn), 1, 0)
		C.XTestFakeButtonEvent(state.XDisplay, C.uint(btn), 0, 0)
	}

	C.XFlush(state.XDisplay)
}

// -------------------- UTILS --------------------

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
