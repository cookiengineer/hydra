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

func QueryAllWindows(bridge *Bridge) ([]types.Window, error) {

	if bridge.display == nil {
		return nil, errors.New("Display is nil")
	}

	var returned_root    C.Window
	var returned_parent  C.Window
	var children         *C.Window
	var num_children     C.uint

	if C.XQueryTree(
		bridge.display,
		bridge.window,
		&returned_root,
		&returned_parent,
		&children,
		&num_children,
	) == 0 {
		return nil, errors.New("XQueryTree failed")
	}

	defer C.XFree(unsafe.Pointer(children))

	c_count := int(num_children)

	go_children := (*[1 << 24]C.Window)(unsafe.Pointer(children))[:c_count:c_count]

	var name_c *C.char
	name_c = C.CString("_NET_WM_NAME")
	net_wm_name := C.XInternAtom(bridge.display, name_c, C.True)
	C.free(unsafe.Pointer(name_c))

	name_c = C.CString("UTF8_STRING")
	utf8_string := C.XInternAtom(bridge.display, name_c, C.True)
	C.free(unsafe.Pointer(name_c))

	windows := make([]types.Window, 0)

	for _, child_window := range go_children {

		var attr C.XWindowAttributes
		if C.XGetWindowAttributes(bridge.display, child_window, &attr) == 0 {
			continue
		}

		if attr.map_state != C.IsViewable {
			continue
		}

		if attr.override_redirect != 0 {
			continue
		}

		var child C.Window
		var abs_x C.int
		var abs_y C.int

		C.XTranslateCoordinates(
			bridge.display,
			child_window,
			bridge.window,
			0, 0,
			&abs_x, &abs_y,
			&child,
		)

		title := ""

		var actual_type C.Atom
		var actual_format C.int
		var num_items C.ulong
		var bytes_after C.ulong
		var data *C.uchar

		if C.XGetWindowProperty(
			bridge.display,
			child_window,
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

		windows = append(windows, types.Window{
			ID:     uint64(child_window),
			Title:  title,
			X:      int(abs_x),
			Y:      int(abs_y),
			Width:  int(attr.width),
			Height: int(attr.height),
		})

	}

	return windows, nil

}
