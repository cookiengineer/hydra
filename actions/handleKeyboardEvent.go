package actions

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func handleKeyboardEvent(bridge *xorg.Bridge, event types.KeyboardEvent, state *types.GlobalState, config *types.Config) {

	if handleWindowManagerAction(bridge, &event, state, config) {
		return
	}

	active := state.GetActive()

	if active == nil {
		return
	}

	if event.Keysym == 0 {

		keysym, err := xorg.KeycodeToKeysym(bridge, event.Keycode)

		if err == nil {
			event.Keysym = keysym
		}

	}

	modifiers, err := bridge.QueryModifiers()

	if err == nil {
		event.Modifiers = modifiers
	}

	if active.Socket != nil {
		sendEnvelope(active, "keyboard", event)
	} else {
		state.ResetActive()
	}

}
