package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
*/
import "C"

import "errors"

func UnfocusWindow(bridge *Bridge) error {

	if bridge.display == nil {
		return errors.New("Display is nil")
	}

	C.XSetInputFocus(
		bridge.display,
		C.None,
		C.RevertToParent,
		C.CurrentTime,
	)

	C.XFlush(bridge.display)

	return nil

}
