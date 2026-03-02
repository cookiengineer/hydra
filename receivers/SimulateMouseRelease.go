package receivers

// TODO: Import C libraries

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateMouseRelease(state *types.State, button int) error {

	if state.XDisplay != nil {

		C.XTestFakeButtonEvent(state.XDisplay, C.uint(button), 0, 0)
		C.XFlush(state.XDisplay)

		return nil

	} else {
		return errors.New("XDisplay is nil")
	}

}
