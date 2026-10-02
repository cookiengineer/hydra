package actions

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func handleMouseEvent(bridge *xorg.Bridge, event types.MouseEvent, state *types.GlobalState, config *types.Config) {

	if config.This != config.Controller {
		return
	}

	mouse_x, mouse_y, err := bridge.QueryPointer()

	if err != nil {
		return
	}

	active := state.GetActive()

	if active == nil {

		controller := config.GetMachine(config.Controller)

		if controller == nil || controller.Screen == nil {
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

		if target != nil && target.Position != "center" && target.Socket != nil {
			activateMouse(bridge, state, config, target, mouse_x, mouse_y)
		}

		return

	}

	virtual_screen := config.GetVirtualScreen()

	if virtual_screen == nil {
		state.ResetActive()
		return
	}

	screen := virtual_screen.GetMachine(active.Hostname)

	if screen == nil {
		state.ResetActive()
		return
	}

	vx, vy := state.GetCursor()

	if event.Type == types.MouseMove {

		vx += event.DX
		vy += event.DY

		min_x := int(screen.OffsetX)
		min_y := int(screen.OffsetY)
		max_x := min_x + int(screen.Width) - 1
		max_y := min_y + int(screen.Height) - 1

		if vx < min_x || vx > max_x || vy < min_y || vy > max_y {
			deactivateRemote(bridge, state, config, active, vx, vy)
			return
		}

		state.SetCursor(vx, vy)

		event.X = uint(vx)
		event.Y = uint(vy)
		event.DX = 0
		event.DY = 0

		sendEnvelope(active, "mouse", event)

		return

	}

	event.X = uint(vx)
	event.Y = uint(vy)

	sendEnvelope(active, "mouse", event)

}
