package xorg

import "github.com/cookiengineer/hydra/types"

func SimulateMouseEvent(bridge *Bridge, event *types.MouseEvent) {

	switch event.Type {
	case types.MouseMove:
		SimulateMouseMove(bridge, event.DX, event.DY)
	case types.MouseButtonPress:
		SimulateMousePress(bridge, event.Button)
	case types.MouseButtonRelease:
		SimulateMouseRelease(bridge, event.Button)
	case types.MouseScroll:
		SimulateMouseScroll(bridge, event.DX, event.DY)
	}

}
