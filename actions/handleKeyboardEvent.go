package actions

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func handleKeyboardEvent(bridge *xorg.Bridge, event types.KeyboardEvent, state *types.GlobalState, config *types.Config) {

	active := state.GetActive()

	if event.Type == types.KeyPress {

		if handleWindowManagerAction(bridge, &event, state, config) {
			return
		}

	}

	if active == nil {
		return
	}

	if event.Keysym == 0 {

		keysym, err := xorg.KeycodeToKeysym(bridge, event.Keycode)

		if err == nil {
			event.Keysym = keysym
		}

	}

	if event.Type == types.KeyPress {

		state.TrackKey(event.Keysym)

	} else {

		if !state.IsTrackedKey(event.Keysym) {
			return
		}

		state.UntrackKey(event.Keysym)

	}

	modifiers, err := bridge.QueryModifiers()

	if err == nil {
		event.Modifiers = modifiers
	}

	if active.Socket != nil {
		sendKeyboard(active, event)
	} else {
		state.ClearTrackedKeys()
		state.ResetActive()
	}

}
