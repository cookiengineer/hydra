package actions

import "encoding/json"
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

	target := config.GetMachine(active.Hostname)

	if target != nil && target.Socket != nil {
		evJSON, err := json.Marshal(event)
		if err == nil {
			select {
			case target.Socket <- evJSON:
			default:
			}
		}
	} else {
		state.ResetActive()
	}

}
