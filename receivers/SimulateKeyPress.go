package receivers

import "errors"
import "github.com/cookiengineer/hydra/types"

func SimulateKeyPress(state *types.State, keycode int) error {

	if state.XDisplay != nil {

		C.XTestFakeKeyEvent(state.XDisplay, C.uint(keycode), 1, 0)
		C.XFlush(state.XDisplay)

		return nil

	} else {
		return errors.New("XDisplay is nil")
	}

}
