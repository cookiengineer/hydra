package helpers

import "github.com/cookiengineer/hydra/types"
import "math"
import "sort"

func FindClosestWindowUp(focused *types.Window, windows []types.Window) *types.Window {

	focused_center_x := focused.X + focused.Width/2

	var candidates []types.Window

	for _, window := range windows {

		if window.ID == focused.ID {
			continue
		}

		w_center_x := window.X + window.Width/2
		x_overlap := math.Abs(float64(focused_center_x - w_center_x))

		if window.Y+window.Height <= focused.Y && x_overlap < float64(focused.Width)/2 {
			candidates = append(candidates, window)
		}

	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		dist_i := focused.Y - (candidates[i].Y + candidates[i].Height)
		dist_j := focused.Y - (candidates[j].Y + candidates[j].Height)
		return dist_i < dist_j
	})

	return &candidates[0]

}
