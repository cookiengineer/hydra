package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
*/
import "C"

import "github.com/cookiengineer/hydra/types"

func handleKeyboardEvent(state *State, cookie *C.XGenericEventCookie) {

	switch cookie.evtype {

	case C.XI_RawKeyPress:

		raw := (*C.XIRawEvent)(cookie.data)

		state.KeyboardEvents <- types.KeyboardEvent{
			Type:    types.KeyPress,
			Keycode: uint32(raw.detail),
		}

	case C.XI_RawKeyRelease:

		raw := (*C.XIRawEvent)(cookie.data)

		state.KeyboardEvents <- types.KeyboardEvent{
			Type:    types.KeyRelease,
			Keycode: uint32(raw.detail),
		}

	}

}

