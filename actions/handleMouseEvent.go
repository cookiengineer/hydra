package actions

import "encoding/json"
import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func handleMouseEvent(bridge *xorg.Bridge, event types.MouseEvent, state *types.GlobalState, config *types.Config) {

	mouse_x, mouse_y, err := bridge.QueryPointer()

	if err != nil {
		return
	}

	active := state.GetActive()

	if config.This != config.Controller {
		return
	}

	if active == nil {

		controller := config.GetMachine(config.Controller)
		if controller == nil {
			return
		}

		var target *types.Machine

		if mouse_x <= 1 {
			target = config.QueryMachine("left-of")
		} else if mouse_x >= int(controller.Screen.Width)-1 {
			target = config.QueryMachine("right-of")
		} else if mouse_y <= 1 {
			target = config.QueryMachine("above")
		} else if mouse_y >= int(controller.Screen.Height)-1 {
			target = config.QueryMachine("below")
		}

		if target != nil && target.Position != "center" {

			focused, err := xorg.QueryFocusedWindow(bridge)

			if err == nil && focused != nil {
				state.SetLastFocusedWindow(focused.ID)
			}

			state.SetActive(target)
			fmt.Printf("Activated remote machine: %s (%s)\n", target.Hostname, target.Position)

			if target.Position == "left-of" {
				bridge.WarpPointer(1, mouse_y)
			} else if target.Position == "right-of" {
				bridge.WarpPointer(int(controller.Screen.Width)-2, mouse_y)
			} else if target.Position == "above" {
				bridge.WarpPointer(mouse_x, 1)
			} else if target.Position == "below" {
				bridge.WarpPointer(mouse_x, int(controller.Screen.Height)-2)
			}
		}

	} else {

		target := config.GetMachine(active.Hostname)

		if target != nil && target.Socket != nil {
			evJSON, err := json.Marshal(event)
			if err == nil {
				select {
				case target.Socket <- evJSON:
				default:
				}
			}
		} else {
			state.ResetActive()
		}

	}

}
