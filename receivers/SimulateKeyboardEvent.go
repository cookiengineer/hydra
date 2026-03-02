package receivers

import "github.com/cookiengineer/hydra/types"

func SimulateKeyboardEvent(state *types.State, event *types.KeyboardEvent) {

	switch event.Type {
	case types.KeyPress:
		SimulateKeyPress(state, event.Key)
	case types.KeyRelease:
		SimulateKeyRelease(state, event.Key)
	}

}
