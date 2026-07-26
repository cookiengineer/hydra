package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
*/
import "C"

import "errors"

func MapWindow(bridge *Bridge, window_id uint64) error {

	if bridge.display == nil {
		return errors.New("Display is nil")
	}

	C.XMapWindow(bridge.display, C.Window(window_id))

	C.XFlush(bridge.display)

	return nil

}
