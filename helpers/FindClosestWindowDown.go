package helpers

import "github.com/cookiengineer/hydra/types"
import "math"
import "sort"

func FindClosestWindowDown(focused *types.Window, windows []types.Window) *types.Window {

	focused_center_x := focused.X + focused.Width/2

	var candidates []types.Window

	for _, window := range windows {

		if window.ID == focused.ID {
			continue
		}

		w_center_x := window.X + window.Width/2
		x_overlap := math.Abs(float64(focused_center_x - w_center_x))

		if window.Y >= focused.Y+focused.Height && x_overlap < float64(focused.Width)/2 {
			candidates = append(candidates, window)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		dist_i := candidates[i].Y - (focused.Y + focused.Height)
		dist_j := candidates[j].Y - (focused.Y + focused.Height)
		return dist_i < dist_j
	})

	return &candidates[0]

}
