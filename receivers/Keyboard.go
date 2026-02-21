package receivers

import "github.com/cookiengineer/hydra/listeners"
import "github.com/cookiengineer/hydra/types"

// ApplyKeyboardEvent applies a keyboard event to the local system
func ApplyKeyboardEvent(state *listeners.State, ke *types.KeyboardEvent) {
	switch ke.Type {
	case types.KeyPress:
		state.SimulateKeyPress(ke.Key)
	case types.KeyRelease:
		state.SimulateKeyRelease(ke.Key)
	}
}
