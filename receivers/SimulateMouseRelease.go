package receivers

// TODO: Import C libraries

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateMouseRelease(state *types.State, button int) error {

	if state.Display != nil {

		C.XTestFakeButtonEvent(state.Display, C.uint(button), 0, 0)
		C.XFlush(state.Display)

		return nil

	} else {
		return errors.New("Display is nil")
	}

}
