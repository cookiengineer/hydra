package receivers

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func ApplyTileEvent(bridge *xorg.Bridge, event *types.TileEvent, virtualScreen *types.VirtualScreen, hostname string) {

	if bridge == nil {
		return
	}

	tile_position := types.TilePositionFromString(event.Position)

	if tile_position == -1 {
		fmt.Printf("ApplyTileEvent: Unknown position %s\n", event.Position)
		return
	}

	window, err := xorg.QueryFocusedWindow(bridge)

	if err != nil || window == nil {
		fmt.Printf("ApplyTileEvent: No focused window\n")
		return
	}

	var monitor *types.Monitor

	local_screen := virtualScreen.GetMachine(hostname)

	if local_screen != nil && len(local_screen.Monitors) > 0 {

		for i := range local_screen.Monitors {

			m := &local_screen.Monitors[i]

			if window.X >= m.OffsetX && window.X < m.OffsetX+m.Width &&
				window.Y >= m.OffsetY && window.Y < m.OffsetY+m.Height {
				monitor = m
				break
			}

		}

		if monitor == nil {
			monitor = &local_screen.Monitors[0]
		}

	}

	xorg.TileWindow(bridge, window.ID, tile_position, monitor)

}
