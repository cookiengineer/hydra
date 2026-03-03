package receivers

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XTest.h>
*/
import "C"

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateKeyRelease(state *types.State, keycode uint32) error {

	if state.Display != nil {

		C.XTestFakeKeyEvent(state.Display, C.uint(keycode), 0, 0)
		C.XFlush(state.Display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
