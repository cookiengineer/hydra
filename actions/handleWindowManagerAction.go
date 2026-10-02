package actions

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/helpers"
import "github.com/cookiengineer/hydra/types"

func handleWindowManagerAction(bridge *xorg.Bridge, event *types.KeyboardEvent, state *types.GlobalState, config *types.Config) bool {

	if event.Type != types.KeyPress {
		return false
	}

	bindings := config.KeyBindings

	modifiers, err := bridge.QueryModifiers()

	if err != nil {
		return false
	}

	for _, binding := range bindings {

		binding_keycode, err := xorg.KeysymToKeycode(bridge, binding.Keycode)

		if err != nil || binding_keycode == 0 {
			continue
		}

		if binding.Modifiers == modifiers && binding_keycode == event.Keycode {

			active := state.GetActive()

			switch binding.Action {

			case types.ActionResetToController:

				if active != nil && active.Socket != nil {
					sendEnvelope(active, "reset", types.ResetEvent{Type: "reset"})
				}

				release_cx, release_cy := state.GetCursor()
				releaseTrackedInput(state, active, release_cx, release_cy)

				state.ResetActive()

				if config.Screen != nil {
					center_x := int(config.Screen.Width) / 2
					center_y := int(config.Screen.Height) / 2
					bridge.WarpPointer(center_x, center_y)
				}

				last_id := state.GetLastFocusedWindow()

				if last_id != 0 {
					xorg.FocusWindow(bridge, last_id)
				}

				fmt.Printf("WindowManager: Reset to controller\n")

				return true

			case types.ActionFocusLeft, types.ActionFocusRight, types.ActionFocusUp, types.ActionFocusDown:

				direction := ""

				switch binding.Action {
				case types.ActionFocusLeft:
					direction = "left"
				case types.ActionFocusRight:
					direction = "right"
				case types.ActionFocusUp:
					direction = "up"
				case types.ActionFocusDown:
					direction = "down"
				}

				if active != nil && active.Socket != nil {

					sendEnvelope(active, "focus", types.FocusEvent{Type: "focus", Direction: direction})

					return true

				}

				windows, err := xorg.QueryAllWindows(bridge)

				if err != nil {
					fmt.Printf("WindowManager: No windows found\n")
					return false
				}

				focused, err := xorg.QueryFocusedWindow(bridge)

				if err != nil || focused == nil {
					fmt.Printf("WindowManager: No focused window\n")
					return false
				}

				var next *types.Window

				switch binding.Action {
				case types.ActionFocusLeft:
					next = helpers.FindClosestWindowLeft(focused, windows)
				case types.ActionFocusRight:
					next = helpers.FindClosestWindowRight(focused, windows)
				case types.ActionFocusUp:
					next = helpers.FindClosestWindowUp(focused, windows)
				case types.ActionFocusDown:
					next = helpers.FindClosestWindowDown(focused, windows)
				}

				if next != nil && next.ID != focused.ID {

					err := xorg.FocusWindow(bridge, next.ID)

					if err == nil {
						fmt.Printf("WindowManager: Focused window %s\n", next.Title)
						return true
					}

				}

				if direction == "right" || direction == "left" {

					position := direction + "-of"
					target := config.QueryMachine(position)

					if target != nil && target.Socket != nil {

						edge := "first"

						if direction == "left" {
							edge = "last"
						}

						return activateFocus(bridge, state, config, target, edge)

					}

					if len(windows) > 0 {

						var wrap *types.Window

						if direction == "right" {
							wrap = helpers.FindLeftmostWindow(windows)
						} else {
							wrap = helpers.FindRightmostWindow(windows)
						}

						if wrap != nil && wrap.ID != focused.ID {
							xorg.FocusWindow(bridge, wrap.ID)
							fmt.Printf("WindowManager: Wrapped to window %s\n", wrap.Title)
							return true
						}

					}

				}

				if direction == "up" || direction == "down" {

					if len(windows) > 0 {

						var wrap *types.Window

						if direction == "up" {
							wrap = helpers.FindBottommostWindow(windows)
						} else {
							wrap = helpers.FindTopmostWindow(windows)
						}

						if wrap != nil && wrap.ID != focused.ID {
							xorg.FocusWindow(bridge, wrap.ID)
							fmt.Printf("WindowManager: Wrapped to window %s\n", wrap.Title)
							return true
						}

					}

				}

				fmt.Printf("WindowManager: No window in that direction\n")

				return false

			case types.ActionTileLeft, types.ActionTileRight, types.ActionTileTop, types.ActionTileBottom, types.ActionTileTopLeft, types.ActionTileTopRight, types.ActionTileBottomLeft, types.ActionTileBottomRight:

				if active != nil && active.Socket != nil {

					tile_position := binding.Action[len("tile-"):]

					sendEnvelope(active, "tile", types.TileEvent{Type: "tile", Position: tile_position})

					return true

				}

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

			case types.ActionSwitchWorkspace:

				ws := state.GetWorkspaceByIndex(binding.Data)

				if ws == nil {
					fmt.Printf("WindowManager: Unknown workspace index %d\n", binding.Data)
					return false
				}

				if active != nil && active.Socket != nil {

					sendEnvelope(active, "workspace", types.WorkspaceEvent{Type: "workspace", Name: ws.Name, Index: ws.Index})

					return true

				}

				helpers.SwitchWorkspace(bridge, state, ws.Name)

				return true

			case types.ActionMoveToWorkspace:

				ws := state.GetWorkspaceByIndex(binding.Data)

				if ws == nil {
					fmt.Printf("WindowManager: Unknown workspace index %d\n", binding.Data)
					return false
				}

				if active != nil && active.Socket != nil {

					sendEnvelope(active, "workspace", types.WorkspaceEvent{Type: "workspace", Name: ws.Name, Index: ws.Index})

					return true

				}

				helpers.MoveWindowToWorkspace(bridge, state, ws.Name)

				return true

			default:

				fmt.Printf("WindowManager: Unknown action %s\n", binding.Action)
				return false

			}

		}

	}

	return false

}
