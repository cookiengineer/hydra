package actions

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func handleWindowManagerAction(bridge *xorg.Bridge, event *types.KeyboardEvent, state *types.GlobalState, config *types.Config) bool {

	if event.Type != types.KeyPress {
		return false
	}

	bindings := types.GetDefaultKeyBindings()

	modifiers, err := bridge.QueryModifiers()

	if err != nil {
		return false
	}

	for _, binding := range bindings {

		if binding.Matches(modifiers, event.Keycode) {

			switch binding.Action {

			case types.ActionResetToController:

				state.ResetActive()

				if config.Screen != nil {
					center_x := int(config.Screen.Width) / 2
					center_y := int(config.Screen.Height) / 2
					bridge.WarpPointer(center_x, center_y)
				}

				fmt.Printf("WindowManager: Reset to controller\n")
				return true

			case types.ActionTileLeft, types.ActionTileRight, types.ActionTileTopLeft, types.ActionTileTopRight, types.ActionTileBottomLeft, types.ActionTileBottomRight:

				tile_position := types.TilePositionFromString(binding.Action[len("tile-"):])

				window, err := xorg.QueryFocusedWindow(bridge)

				if err == nil && window != nil {

					var monitor *types.Monitor

					controller_screen := config.Screen.GetMachine(config.Controller)

					if controller_screen != nil && len(controller_screen.Monitors) > 0 {
						for i := range controller_screen.Monitors {
							m := &controller_screen.Monitors[i]

							if window.X >= m.OffsetX && window.X < m.OffsetX+m.Width &&
								window.Y >= m.OffsetY && window.Y < m.OffsetY+m.Height {
								monitor = m
								break
							}
						}

						if monitor == nil {
							monitor = &controller_screen.Monitors[0]
						}
					}

					err1 := xorg.TileWindow(bridge, window.ID, tile_position, monitor)

					if err1 == nil {
						fmt.Printf("WindowManager: Tiled window to %s\n", tile_position.String())
						return true
					} else {
						fmt.Printf("WindowManager: Error tiling window: %s\n", err1.Error())
						return false
					}

				} else {
					fmt.Printf("WindowManager: No focused window found\n")
					return false
				}

			default:
				fmt.Printf("WindowManager: Unknown action %s\n", binding.Action)
				return false
			}
		}
	}

	return false
}
