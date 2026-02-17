package listeners

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>
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

func Init(display string) (*State, error) {

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

	return &State{
		XDisplay:       x_display,
		XIOpcode:       xi_opcode,
		XWindow:        x_root_window,
		MouseEvents:    make(chan types.MouseEvent, 32),
		KeyboardEvents: make(chan types.KeyboardEvent, 32),
	}, nil

}

