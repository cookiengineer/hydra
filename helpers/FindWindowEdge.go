package helpers

import "github.com/cookiengineer/hydra/types"

func FindLeftmostWindow(windows []types.Window) *types.Window {

	if len(windows) == 0 {
		return nil
	}

	best := &windows[0]

	for i := range windows {

		if windows[i].X < best.X || (windows[i].X == best.X && windows[i].Y < best.Y) {
			best = &windows[i]
		}

	}

	return best

}

func FindRightmostWindow(windows []types.Window) *types.Window {

	if len(windows) == 0 {
		return nil
	}

	best := &windows[0]

	for i := range windows {

		if windows[i].X+windows[i].Width > best.X+best.Width {
			best = &windows[i]
		}

	}

	return best

}

func FindTopmostWindow(windows []types.Window) *types.Window {

	if len(windows) == 0 {
		return nil
	}

	best := &windows[0]

	for i := range windows {

		if windows[i].Y < best.Y {
			best = &windows[i]
		}

	}

	return best

}

func FindBottommostWindow(windows []types.Window) *types.Window {

	if len(windows) == 0 {
		return nil
	}

	best := &windows[0]

	for i := range windows {

		if windows[i].Y+windows[i].Height > best.Y+best.Height {
			best = &windows[i]
		}

	}

	return best

}
