package receivers

// TODO: Import C libraries

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateMousePress(state *types.State, button types.MouseEventButton) error {

	if state.XDisplay != nil {

		C.XTestFakeButtonEvent(state.XDisplay, C.uint(uint(button)), 1, 0)
		C.XFlush(state.XDisplay)

		return nil

	} else {
		return errors.New("XDisplay is nil")
	}

}
