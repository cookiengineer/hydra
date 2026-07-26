package receivers

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func ApplyKeyboardEvent(bridge *xorg.Bridge, event *types.KeyboardEvent) {

	if bridge != nil {
		xorg.SimulateKeyboardEvent(bridge, event)
	}

}
