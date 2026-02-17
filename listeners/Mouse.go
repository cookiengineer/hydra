package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
*/
import "C"

import "math"
import "unsafe"
import "github.com/cookiengineer/hydra/types"

func handleMouseEvent(state *State, cookie *C.XGenericEventCookie) {

	switch cookie.evtype {

	case C.XI_RawMotion:

		raw := (*C.XIRawEvent)(cookie.data)

		if raw.valuators.mask_len >= 2 {

			values := (*[2]float64)(unsafe.Pointer(raw.raw_values))

			if math.Abs(values[0]) < 0.01 || math.Abs(values[1]) < 0.01 {

				// Ignore Xinput fake motion event which has dx set to scroll distance

			} else {
				state.MouseEvents <- types.MouseEvent{
					Type: types.MouseMove,
					DX:   values[0],
					DY:   values[1],
				}
			}

		}

	case C.XI_RawButtonPress:

		raw := (*C.XIRawEvent)(cookie.data)

		switch int(raw.detail) {
		// 1 = Left Button
		case 1:
			state.MouseEvents <- types.MouseEvent{
				Type:   types.MouseButtonPress,
				DX:     0,
				DY:     0,
				Button: types.MouseButtonLeft,
			}
		// 2 = Middle Button
		case 2:
			state.MouseEvents <- types.MouseEvent{
				Type:   types.MouseButtonPress,
				DX:     0,
				DY:     0,
				Button: types.MouseButtonMiddle,
			}
		// 3 = Right Button
		case 3:
			state.MouseEvents <- types.MouseEvent{
				Type:   types.MouseButtonPress,
				DX:     0,
				DY:     0,
				Button: types.MouseButtonRight,
			}
		// 4 = Scroll Wheel Up
		case 4:
			state.MouseEvents <- types.MouseEvent{
				Type: types.MouseScroll,
				DX:   0,
				DY:   1,
			}
		// 5 = Scroll Wheel Down
		case 5:
			state.MouseEvents <- types.MouseEvent{
				Type: types.MouseScroll,
				DX:   0,
				DY:  -1,
			}
		// 6 = Scroll Wheel Left
		case 6:
			state.MouseEvents <- types.MouseEvent{
				Type: types.MouseScroll,
				DX:   -1,
				DY:    0,
			}
		// 7 = Scroll Wheel Right
		case 7:
			state.MouseEvents <- types.MouseEvent{
				Type: types.MouseScroll,
				DX:   1,
				DY:   0,
			}
		default:
			// XXX: Debug info
			// state.MouseEvents <- types.MouseEvent{
			// 	Type:   types.MouseButtonPress,
			// 	Button: int(raw.detail),
			// }
		}

	case C.XI_RawButtonRelease:

		raw := (*C.XIRawEvent)(cookie.data)

		switch int(raw.detail) {
		// 1 = Left Button
		case 1:
			state.MouseEvents <- types.MouseEvent{
				Type:   types.MouseButtonRelease,
				DX:     0,
				DY:     0,
				Button: types.MouseButtonLeft,
			}
		// 2 = Middle Button
		case 2:
			state.MouseEvents <- types.MouseEvent{
				Type:   types.MouseButtonRelease,
				DX:     0,
				DY:     0,
				Button: types.MouseButtonMiddle,
			}
		// 3 = Right Button
		case 3:
			state.MouseEvents <- types.MouseEvent{
				Type:   types.MouseButtonRelease,
				DX:     0,
				DY:     0,
				Button: types.MouseButtonRight,
			}

		}

	}

}

