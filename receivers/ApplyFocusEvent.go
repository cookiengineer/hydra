package receivers

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/helpers"
import "github.com/cookiengineer/hydra/types"

func ApplyFocusEvent(bridge *xorg.Bridge, event *types.FocusEvent) {

	if bridge == nil {
		return
	}

	windows, err := xorg.QueryAllWindows(bridge)

	if err != nil || len(windows) == 0 {
		fmt.Printf("ApplyFocusEvent: No windows found\n")
		return
	}

	focused, err := xorg.QueryFocusedWindow(bridge)

	if err != nil || focused == nil {
		fmt.Printf("ApplyFocusEvent: No focused window\n")
		return
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
		return
	}

	if next != nil && next.ID != focused.ID {
		xorg.FocusWindow(bridge, next.ID)
	}

}
