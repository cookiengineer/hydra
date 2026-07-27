//go:build integration

package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdlib.h>
#include <string.h>

static Window create_test_window(Display *display, Window root, int x, int y, int w, int h, const char *title) {

	Window win = XCreateSimpleWindow(
		display,
		root,
		x, y, w, h,
		0,
		BlackPixel(display, DefaultScreen(display)),
		WhitePixel(display, DefaultScreen(display))
	);

	if (title != NULL) {
		Atom net_wm_name = XInternAtom(display, "_NET_WM_NAME", False);
		Atom utf8_string = XInternAtom(display, "UTF8_STRING", False);
		XChangeProperty(
			display, win,
			net_wm_name, utf8_string, 8,
			PropModeReplace,
			(const unsigned char *)title,
			(int)strlen(title)
		);
	}

	XMapWindow(display, win);
	XSync(display, False);

	return win;
}
*/
import "C"

import "errors"
import "unsafe"

func createTestWindow(bridge *Bridge, title string, x int, y int, w int, h int) (uint64, error) {

	if bridge.display == nil {
		return 0, errors.New("XDisplay is nil")
	}

	c_title := C.CString(title)
	defer C.free(unsafe.Pointer(c_title))

	window := C.create_test_window(
		bridge.display,
		bridge.window,
		C.int(x),
		C.int(y),
		C.int(w),
		C.int(h),
		c_title,
	)

	if window == 0 {
		return 0, errors.New("XCreateSimpleWindow failed")
	}

	return uint64(window), nil

}

