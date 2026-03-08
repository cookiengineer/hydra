package xorg

import "github.com/cookiengineer/hydra/types"

func SimulateMouseEvent(bridge *Bridge, event *types.MouseEvent) {

	switch event.Type {
	case types.MouseMove:
		SimulateMouseMove(bridge, event.X, event.Y, event.DX, event.DY)
	case types.MouseButtonPress:
		SimulateMousePress(bridge, event.X, event.Y, event.Button)
	case types.MouseButtonRelease:
		SimulateMouseRelease(bridge, event.X, event.Y, event.Button)
	case types.MouseScroll:
		SimulateMouseScroll(bridge, event.X, event.Y, event.DX, event.DY)
	}

}
