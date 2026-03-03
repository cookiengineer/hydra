package receivers

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

func SimulateMouseRelease(state *types.State, button types.MouseEventButton) error {

	if state.Display != nil {

		C.XTestFakeButtonEvent(state.Display, C.uint(uint(button)), 0, 0)
		C.XFlush(state.Display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
