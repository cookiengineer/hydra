package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XInput2.h>
#include <string.h>
#include <stdlib.h>

static int init_xinput(Display *display) {
	int major = 2, minor = 0;
	return XIQueryVersion(display, &major, &minor);
}

static void select_events(Display *display, Window win) {
	XIEventMask evmask;
	unsigned char mask[(XI_LASTEVENT + 7)/8];
	memset(mask, 0, sizeof(mask));

	XISetMask(mask, XI_RawMotion);
	XISetMask(mask, XI_RawButtonPress);
	XISetMask(mask, XI_RawButtonRelease);

	evmask.deviceid = XIAllMasterDevices;
	evmask.mask_len = sizeof(mask);
	evmask.mask = mask;

	XISelectEvents(display, win, &evmask, 1);
	XFlush(display);
}
*/
import "C"

import "github.com/cookiengineer/hydra/types"
import "errors"
import "math"
import "os"
import "unsafe"

func CaptureMouse(events chan<- types.MouseEvent, xdisplay string) error {

	if xdisplay != "" {
		os.Setenv("DISPLAY", xdisplay)
	}

	display := C.XOpenDisplay(nil)

	if display != nil {

		root := C.XDefaultRootWindow(display)

		if C.init_xinput(display) == C.Success {

			C.select_events(display, root)

			go func() {

				defer C.XCloseDisplay(display)

				for {

					var event C.XEvent

					C.XNextEvent(display, &event)

					// I am starting to hate unions
					event_type := *(*C.int)(unsafe.Pointer(&event))

					if event_type == C.GenericEvent {

						// Handle potential XInput event now

						cookie := (*C.XGenericEventCookie)(unsafe.Pointer(&event))

						if C.XGetEventData(display, cookie) == 0 {
							continue
						}

						switch cookie.evtype {

						case C.XI_RawMotion:

							raw := (*C.XIRawEvent)(cookie.data)

							if raw.valuators.mask_len >= 2 {

								values := (*[2]float64)(unsafe.Pointer(raw.raw_values))

								if math.Abs(values[0]) < 0.01 || math.Abs(values[1]) < 0.01 {

									// Ignore Xinput fake motion event which has dx set to scroll distance

								} else {
									events <- types.MouseEvent{
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
								events <- types.MouseEvent{
									Type:   types.MouseButtonPress,
									DX:     0,
									DY:     0,
									Button: types.MouseButtonLeft,
								}
							// 2 = Middle Button
							case 2:
								events <- types.MouseEvent{
									Type:   types.MouseButtonPress,
									DX:     0,
									DY:     0,
									Button: types.MouseButtonMiddle,
								}
							// 3 = Right Button
							case 3:
								events <- types.MouseEvent{
									Type:   types.MouseButtonPress,
									DX:     0,
									DY:     0,
									Button: types.MouseButtonRight,
								}
							// 4 = Scroll Wheel Up
							case 4:
								events <- types.MouseEvent{
									Type: types.MouseScroll,
									DX:   0,
									DY:   1,
								}
							// 5 = Scroll Wheel Down
							case 5:
								events <- types.MouseEvent{
									Type: types.MouseScroll,
									DX:   0,
									DY:  -1,
								}
							// 6 = Scroll Wheel Left
							case 6:
								events <- types.MouseEvent{
									Type: types.MouseScroll,
									DX:   -1,
									DY:    0,
								}
							// 7 = Scroll Wheel Right
							case 7:
								events <- types.MouseEvent{
									Type: types.MouseScroll,
									DX:   1,
									DY:   0,
								}
							default:
								// XXX: Debug info
								// events <- types.MouseEvent{
								// 	Type:   types.MouseButtonPress,
								// 	Button: int(raw.detail),
								// }
							}

						case C.XI_RawButtonRelease:

							raw := (*C.XIRawEvent)(cookie.data)

							switch int(raw.detail) {
							// 1 = Left Button
							case 1:
								events <- types.MouseEvent{
									Type:   types.MouseButtonRelease,
									DX:     0,
									DY:     0,
									Button: types.MouseButtonLeft,
								}
							// 2 = Middle Button
							case 2:
								events <- types.MouseEvent{
									Type:   types.MouseButtonRelease,
									DX:     0,
									DY:     0,
									Button: types.MouseButtonMiddle,
								}
							// 3 = Right Button
							case 3:
								events <- types.MouseEvent{
									Type:   types.MouseButtonRelease,
									DX:     0,
									DY:     0,
									Button: types.MouseButtonRight,
								}

							}

						}

						C.XFreeEventData(display, cookie)

					}

				}

			}()

			return nil

		} else {
			return errors.New("Cannot open X default root window via xinput2")
		}

	} else {
		return errors.New("Cannot open X display")
	}

}

