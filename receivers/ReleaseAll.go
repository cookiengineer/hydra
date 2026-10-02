package receivers

import "sync"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

var pressedMutex sync.Mutex
var pressedKeys = make(map[uint32]bool)
var pressedButtons = make(map[types.MouseEventButton]bool)

func trackKey(keycode uint32, down bool) {

	if keycode == 0 {
		return
	}

	pressedMutex.Lock()
	defer pressedMutex.Unlock()

	if down {
		pressedKeys[keycode] = true
	} else {
		delete(pressedKeys, keycode)
	}

}

func trackButton(button types.MouseEventButton, down bool) {

	pressedMutex.Lock()
	defer pressedMutex.Unlock()

	if down {
		pressedButtons[button] = true
	} else {
		delete(pressedButtons, button)
	}

}

func ReleaseAll(bridge *xorg.Bridge) {

	if bridge == nil {
		return
	}

	pressedMutex.Lock()

	keys := make([]uint32, 0, len(pressedKeys))

	for keycode := range pressedKeys {
		keys = append(keys, keycode)
	}

	buttons := make([]types.MouseEventButton, 0, len(pressedButtons))

	for button := range pressedButtons {
		buttons = append(buttons, button)
	}

	pressedKeys = make(map[uint32]bool)
	pressedButtons = make(map[types.MouseEventButton]bool)

	pressedMutex.Unlock()

	for _, keycode := range keys {
		xorg.SimulateKeyRelease(bridge, keycode)
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		x = 0
		y = 0
	}

	for _, button := range buttons {
		xorg.SimulateMouseRelease(bridge, uint(x), uint(y), button)
	}

	xorg.ReleaseModifiers(bridge)

}
