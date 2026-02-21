package receivers

import "github.com/cookiengineer/hydra/listeners"
import "github.com/cookiengineer/hydra/types"

// ApplyMouseEvent applies a mouse event to the local system
func ApplyMouseEvent(state *listeners.State, me *types.MouseEvent) {
	switch me.Type {
	case types.MouseMove:
		state.WarpPointer(me.DX, me.DY)
	case types.MouseButtonPress:
		state.SimulateMousePress(me.Button)
	case types.MouseButtonRelease:
		state.SimulateMouseRelease(me.Button)
	case types.MouseScroll:
		state.SimulateMouseScroll(me.DX, me.DY)
	}
}
