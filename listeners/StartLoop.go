package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
*/
import "C"

import "unsafe"

func StartLoop(state *State) {

	if state != nil {

		if state.running == false {

			state.running = true

			go func() {

				for {

					var event C.XEvent
					C.XNextEvent(state.XDisplay, &event)

					eventType := *(*C.int)(unsafe.Pointer(&event))

					if eventType != C.GenericEvent {
						continue
					}

					cookie := (*C.XGenericEventCookie)(unsafe.Pointer(&event))

					if cookie.extension != state.XIOpcode {
						continue
					}

					if C.XGetEventData(state.XDisplay, cookie) == 0 {
						continue
					}

					switch cookie.evtype {

					case C.XI_RawMotion:
						handleMouseEvent(state, cookie)
					case C.XI_RawButtonPress:
						handleMouseEvent(state, cookie)
					case C.XI_RawButtonRelease:
						handleMouseEvent(state, cookie)
					case C.XI_RawKeyPress:
						handleKeyboardEvent(state, cookie)
					case C.XI_RawKeyRelease:
						handleKeyboardEvent(state, cookie)
					}

					C.XFreeEventData(state.XDisplay, cookie)

				}

			}()

		}

	}

}

