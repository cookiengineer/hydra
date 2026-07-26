package receivers

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func ApplyMouseEvent(bridge *xorg.Bridge, event *types.MouseEvent, virtualScreen *types.VirtualScreen, hostname string) {

	local_screen := virtualScreen.GetMachine(hostname)

	var local_x uint = event.X
	var local_y uint = event.Y

	if local_screen != nil {

		if event.X >= local_screen.OffsetX {
			local_x = event.X - local_screen.OffsetX
		}

		if event.Y >= local_screen.OffsetY {
			local_y = event.Y - local_screen.OffsetY
		}

	}

	event.X = local_x
	event.Y = local_y

	if bridge != nil {
		xorg.SimulateMouseEvent(bridge, event)
	}

}

