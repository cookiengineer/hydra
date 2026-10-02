package receivers

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func ApplyKeyboardEvent(bridge *xorg.Bridge, event *types.KeyboardEvent) {

	if bridge == nil {
		return
	}

	if event.Keysym != 0 {

		keycode, err := xorg.KeysymToKeycode(bridge, event.Keysym)

		if err == nil && keycode != 0 {
			event.Keycode = keycode
		}

	}

	xorg.SimulateKeyboardEvent(bridge, event)

}
