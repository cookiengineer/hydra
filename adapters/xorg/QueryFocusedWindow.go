package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11 -lXi -lXtst
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <stdlib.h>
*/
import "C"

import "errors"
import "unsafe"
import "github.com/cookiengineer/hydra/types"

func QueryFocusedWindow(bridge *Bridge) (*types.Window, error) {

	if bridge.display == nil {
		return nil, errors.New("Display is nil")
	}

	var revert_to C.int
	var focused_window C.Window

	C.XGetInputFocus(bridge.display, &focused_window, &revert_to)

	if focused_window == 0 || focused_window == bridge.window {
		return nil, errors.New("No focused window")
	}

	var attr C.XWindowAttributes
	if C.XGetWindowAttributes(bridge.display, focused_window, &attr) == 0 {
		return nil, errors.New("Cannot get window attributes")
	}

	var child C.Window
	var abs_x C.int
	var abs_y C.int

	C.XTranslateCoordinates(
		bridge.display,
		focused_window,
		bridge.window,
		0, 0,
		&abs_x, &abs_y,
		&child,
	)

	var net_wm_name C.Atom
	name_c := C.CString("_NET_WM_NAME")
	net_wm_name = C.XInternAtom(bridge.display, name_c, C.True)
	C.free(unsafe.Pointer(name_c))

	var utf8_string C.Atom
	utf8_c := C.CString("UTF8_STRING")
	utf8_string = C.XInternAtom(bridge.display, utf8_c, C.True)
	C.free(unsafe.Pointer(utf8_c))

	var actual_type C.Atom
	var actual_format C.int
	var num_items C.ulong
	var bytes_after C.ulong
	var data *C.uchar

	title := ""

	if C.XGetWindowProperty(
		bridge.display,
		focused_window,
		net_wm_name,
		0, 1024,
		C.False,
		utf8_string,
		&actual_type,
		&actual_format,
		&num_items,
		&bytes_after,
		&data,
	) == C.Success && data != nil {
		title = C.GoString((*C.char)(unsafe.Pointer(data)))
		C.XFree(unsafe.Pointer(data))
	}

	if title == "" {
		wm_name_c := C.CString("WM_NAME")
		wm_name := C.XInternAtom(bridge.display, wm_name_c, C.True)
		C.free(unsafe.Pointer(wm_name_c))

		var text_prop C.XTextProperty
		if C.XGetTextProperty(bridge.display, focused_window, &text_prop, wm_name) != 0 {
			if text_prop.nitems > 0 {
				title = C.GoString((*C.char)(unsafe.Pointer(text_prop.value)))
				C.XFree(unsafe.Pointer(text_prop.value))
			}
		}
	}

	return &types.Window{
		ID:     uint64(focused_window),
		Title:  title,
		X:      int(abs_x),
		Y:      int(abs_y),
		Width:  int(attr.width),
		Height: int(attr.height),
	}, nil

}
