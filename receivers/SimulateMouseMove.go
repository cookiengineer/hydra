package receivers

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
*/
import "C"

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateMouseMove(state *types.State, x int, y int) error {

	if state.XDisplay != nil {

		dest_x := C.int(x)
		dest_y := C.int(y)

		if state.VirtualScreen != nil && state.VirtualScreen.Active != nil {
			offset_x := state.VirtualScreen.Active.OffsetX
			offset_y := state.VirtualScreen.Active.OffsetY
			dest_x = C.int(x - offset_x)
			dest_y = C.int(y - offset_y)
		}

		C.XWarpPointer(
			state.XDisplay,
			0,
			state.XWindow,
			0, 0, 0, 0,
			dest_x,
			dest_y,
		)
		C.XFlush(state.XDisplay)

		return nil

	} else {
		return errors.New("XDisplay is nil")
	}

}


