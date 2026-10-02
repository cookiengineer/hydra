package receivers

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/helpers"
import "github.com/cookiengineer/hydra/types"

func ApplyFocusEvent(bridge *xorg.Bridge, event *types.FocusEvent) bool {

	if bridge == nil {
		return false
	}

	windows, err := xorg.QueryAllWindows(bridge)

	if err != nil || len(windows) == 0 {
		fmt.Printf("ApplyFocusEvent: No windows found\n")
		return false
	}

	if event.Edge == "first" {

		next := helpers.FindLeftmostWindow(windows)

		if next != nil {
			xorg.FocusWindow(bridge, next.ID)
			return true
		}

		return false

	}

	if event.Edge == "last" {

		next := helpers.FindRightmostWindow(windows)

		if next != nil {
			xorg.FocusWindow(bridge, next.ID)
			return true
		}

		return false

	}

	focused, err := xorg.QueryFocusedWindow(bridge)

	if err != nil || focused == nil {
		fmt.Printf("ApplyFocusEvent: No focused window\n")
		return false
	}

	var next *types.Window

	switch event.Direction {
	case "left":
		next = helpers.FindClosestWindowLeft(focused, windows)
	case "right":
		next = helpers.FindClosestWindowRight(focused, windows)
	case "up":
		next = helpers.FindClosestWindowUp(focused, windows)
	case "down":
		next = helpers.FindClosestWindowDown(focused, windows)
	default:
		fmt.Printf("ApplyFocusEvent: Unknown direction %s\n", event.Direction)
		return false
	}

	if next != nil && next.ID != focused.ID {

		err := xorg.FocusWindow(bridge, next.ID)

		if err == nil {
			return true
		}

	}

	return false

}
