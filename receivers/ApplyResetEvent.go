package receivers

import "github.com/cookiengineer/hydra/adapters/xorg"

func ApplyResetEvent(bridge *xorg.Bridge) {

	if bridge != nil {
		xorg.UnfocusWindow(bridge)
	}

}
