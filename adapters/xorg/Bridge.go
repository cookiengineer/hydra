package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XInput2.h>
#include <X11/extensions/XTest.h>
#include <stdlib.h>
#include <string.h>

static void register_xinput_events(Display *display, Window win) {

	XIEventMask evmask;
	unsigned char mask[(XI_LASTEVENT + 7)/8];
	memset(mask, 0, sizeof(mask));

	// Mouse
	XISetMask(mask, XI_RawMotion);
	XISetMask(mask, XI_RawButtonPress);
	XISetMask(mask, XI_RawButtonRelease);

	// Keyboard
	XISetMask(mask, XI_RawKeyPress);
	XISetMask(mask, XI_RawKeyRelease);

	evmask.deviceid = XIAllMasterDevices;
	evmask.mask_len = sizeof(mask);
	evmask.mask = mask;

	XISelectEvents(display, win, &evmask, 1);
	XFlush(display);

}
*/
import "C"

import "errors"
import "os"
import "unsafe"
import "github.com/cookiengineer/hydra/types"

type Bridge struct {
	display       *C.Display
	opcode         C.int
	window         C.Window
	MouseEvents    chan types.MouseEvent
	KeyboardEvents chan types.KeyboardEvent
	Screen         *types.Screen
	running        bool
}

func NewBridge(display string) (*Bridge, error) {

	if display != "" {
		os.Setenv("DISPLAY", display)
	}

	x_display := C.XOpenDisplay(nil)

	if x_display == nil {
		return nil, errors.New("Cannot open X display")
	}

	extension := C.CString("XInputExtension")
	defer C.free(unsafe.Pointer(extension))

	var xi_opcode C.int
	var event C.int
	var err C.int

	if C.XQueryExtension(x_display, extension, &xi_opcode, &event, &err) == 0 {

		C.XCloseDisplay(x_display)

		return nil, errors.New("XInput extension not available")

	}

	major := C.int(2)
	minor := C.int(0)

	if C.XIQueryVersion(x_display, &major, &minor) != C.Success {

		C.XCloseDisplay(x_display)

		return nil, errors.New("XInput2 not supported")

	}

	x_root_window := C.XDefaultRootWindow(x_display)

	C.register_xinput_events(x_display, x_root_window)

	return &Bridge{
		display:        x_display,
		opcode:         xi_opcode,
		window:         x_root_window,
		MouseEvents:    make(chan types.MouseEvent, 32),
		KeyboardEvents: make(chan types.KeyboardEvent, 32),
	}, nil

}

func (bridge *Bridge) Destroy() {

	if bridge.display != nil {

		C.XCloseDisplay(bridge.display)
		bridge.display = nil

	}

}

func (bridge *Bridge) Init() {

	if bridge.running == false {

		bridge.running = true

		go func() {

			for {

				var event C.XEvent

				C.XNextEvent(bridge.display, &event)

				eventType := *(*C.int)(unsafe.Pointer(&event))

				if eventType != C.GenericEvent {
					continue
				}

				cookie := (*C.XGenericEventCookie)(unsafe.Pointer(&event))

				if cookie.extension != bridge.opcode {
					continue
				}

				if C.XGetEventData(bridge.display, cookie) == 0 {
					continue
				}

				switch cookie.evtype {

				case C.XI_RawMotion:
					HandleMouseEvent(bridge, cookie)
				case C.XI_RawButtonPress:
					HandleMouseEvent(bridge, cookie)
				case C.XI_RawButtonRelease:
					HandleMouseEvent(bridge, cookie)
				case C.XI_RawKeyPress:
					HandleKeyboardEvent(bridge, cookie)
				case C.XI_RawKeyRelease:
					HandleKeyboardEvent(bridge, cookie)
				}

				C.XFreeEventData(bridge.display, cookie)

			}

		}()

	}

}

func (bridge *Bridge) QueryPointer() (int, int, error) {

	if bridge.display != nil {

		var returned_root     C.Window
		var returned_window   C.Window
		var returned_root_x   C.int
		var returned_root_y   C.int
		var returned_window_x C.int
		var returned_window_y C.int
		var returned_mask     C.uint

		result := C.XQueryPointer(
			bridge.display,
			bridge.window,
			&returned_root,
			&returned_window,
			&returned_root_x,
			&returned_root_y,
			&returned_window_x,
			&returned_window_y,
			&returned_mask,
		)

		if result == 0 {
			return 0, 0, errors.New("XQueryPointer failed")
		} else {
			return int(returned_root_x), int(returned_root_y), nil
		}

	} else {
		return 0, 0, errors.New("XDisplay is nil")
	}

}

