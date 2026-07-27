package xorg

import "github.com/cookiengineer/hydra/types"

func SimulateKeyboardEvent(bridge *Bridge, event *types.KeyboardEvent) {

	if bridge == nil {
		return
	}

	switch event.Type {
	case types.KeyPress:
		SimulateKeyPress(bridge, event.Keycode)
	case types.KeyRelease:
		SimulateKeyRelease(bridge, event.Keycode)
	}

}
