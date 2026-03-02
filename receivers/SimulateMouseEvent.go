package receivers

import "github.com/cookiengineer/hydra/types"

func SimulateMouseEvent(state *types.State, event *types.MouseEvent) {

	switch event.Type {
	case types.MouseMove:
		SimulateMouseMove(state, event.DX, event.DY)
	case types.MouseButtonPress:
		SimulateMousePress(state, event.Button)
	case types.MouseButtonRelease:
		SimulateMouseRelease(state, event.Button)
	case types.MouseScroll:
		SimulateMouseScroll(state, event.DX, event.DY)
	}

}
