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

	if state.XDisplay != nil {

		C.XTestFakeKeyEvent(state.XDisplay, C.uint(keycode), 0, 0)
		C.XFlush(state.XDisplay)

		return nil

	} else {
		return errors.New("XDisplay is nil")
	}

}
