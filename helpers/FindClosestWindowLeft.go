package helpers

import "github.com/cookiengineer/hydra/types"
import "math"
import "sort"

func FindClosestWindowLeft(focused *types.Window, windows []types.Window) *types.Window {

	focused_center_y := focused.Y + focused.Height/2

	var candidates []types.Window

	for _, window := range windows {

		if window.ID == focused.ID {
			continue
		}

		w_center_y := window.Y + window.Height/2
		y_overlap := math.Abs(float64(focused_center_y - w_center_y))

		if window.X+window.Width <= focused.X && y_overlap < float64(focused.Height)/2 {
			candidates = append(candidates, window)
		}

	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		dist_i := focused.X - (candidates[i].X + candidates[i].Width)
		dist_j := focused.X - (candidates[j].X + candidates[j].Width)
		return dist_i < dist_j
	})

	return &candidates[0]

}

