package xorg

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/keysym.h>
#include <X11/XKBlib.h>
*/
import "C"

import "errors"

func KeysymToKeycode(bridge *Bridge, keysym uint32) (uint32, error) {

	if bridge.display == nil {
		return 0, errors.New("Display is nil")
	}

	keycode := C.XKeysymToKeycode(bridge.display, C.KeySym(keysym))

	return uint32(keycode), nil

}

func KeycodeToKeysym(bridge *Bridge, keycode uint32) (uint32, error) {

	if bridge.display == nil {
		return 0, errors.New("Display is nil")
	}

	keysym := C.XkbKeycodeToKeysym(bridge.display, C.KeyCode(keycode), 0, 0)

	return uint32(keysym), nil

}
