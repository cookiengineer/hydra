//go:build integration

package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdlib.h>
#include <string.h>

static void destroy_test_window(Display *display, Window win) {

	XDestroyWindow(display, win);
	XSync(display, False);

}
*/
import "C"

func destroyTestWindow(bridge *Bridge, window_id uint64) {

	if bridge.display == nil {
		return
	}

	C.destroy_test_window(bridge.display, C.Window(window_id))

}
