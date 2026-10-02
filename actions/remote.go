package actions

import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/helpers"
import "github.com/cookiengineer/hydra/types"

func clamp(value int, min int, max int) int {

	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value

}

func sendEnvelope(machine *types.Machine, event_type string, payload interface{}) {

	if machine == nil || machine.Socket == nil {
		return
	}

	envelope, err := types.NewEvent(event_type, payload)

	if err != nil {
		return
	}

	data, err := json.Marshal(envelope)

	if err != nil {
		return
	}

	select {
	case machine.Socket <- data:
	default:
	}

}

func controllerScreen(config *types.Config) *types.Screen {

	if config == nil {
		return nil
	}

	machine := config.GetMachine(config.Controller)

	if machine == nil {
		return nil
	}

	return machine.Screen

}

func activateMouse(bridge *xorg.Bridge, state *types.GlobalState, config *types.Config, target *types.Machine, mouse_x int, mouse_y int) {

	virtual_screen := config.GetVirtualScreen()

	if virtual_screen == nil {
		return
	}

	screen := virtual_screen.GetMachine(target.Hostname)

	if screen == nil {
		return
	}

	focused, err := xorg.QueryFocusedWindow(bridge)

	if err == nil && focused != nil {
		state.SetLastFocusedWindow(focused.ID)
	}

	var vx, vy int

	switch target.Position {

	case "right-of":
		vx = int(screen.OffsetX)
		vy = clamp(mouse_y, int(screen.OffsetY), int(screen.OffsetY)+int(screen.Height)-1)

	case "left-of":
		vx = int(screen.OffsetX) + int(screen.Width) - 1
		vy = clamp(mouse_y, int(screen.OffsetY), int(screen.OffsetY)+int(screen.Height)-1)

	case "above":
		vx = clamp(mouse_x, int(screen.OffsetX), int(screen.OffsetX)+int(screen.Width)-1)
		vy = int(screen.OffsetY) + int(screen.Height) - 1

	case "below":
		vx = clamp(mouse_x, int(screen.OffsetX), int(screen.OffsetX)+int(screen.Width)-1)
		vy = int(screen.OffsetY)

	default:
		return

	}

	state.SetActive(target)
	state.SetCursor(vx, vy)

	xorg.UnfocusWindow(bridge)

	sendEnvelope(target, "mouse", types.MouseEvent{
		Type: types.MouseMove,
		X:    uint(vx),
		Y:    uint(vy),
	})

	fmt.Printf("Activated remote machine: %s (%s)\n", target.Hostname, target.Position)

}

func deactivateRemote(bridge *xorg.Bridge, state *types.GlobalState, config *types.Config, active *types.Machine, vx int, vy int) {

	state.ResetActive()

	screen := controllerScreen(config)

	if screen != nil {

		x := int(screen.Width) - 2
		y := clamp(vy, 1, int(screen.Height)-2)

		switch active.Position {

		case "left-of":
			x = 1

		case "right-of":
			x = int(screen.Width) - 2

		case "above":
			y = 1

		case "below":
			y = int(screen.Height) - 2

		}

		if active.Position == "left-of" || active.Position == "right-of" {
			bridge.WarpPointer(x, y)
		} else {
			bridge.WarpPointer(clamp(vx, 1, int(screen.Width)-2), y)
		}

	}

	last_id := state.GetLastFocusedWindow()

	if last_id != 0 {
		xorg.FocusWindow(bridge, last_id)
	}

	fmt.Printf("Deactivated remote machine: %s\n", active.Hostname)

}

func activateFocus(bridge *xorg.Bridge, state *types.GlobalState, config *types.Config, target *types.Machine, edge string) bool {

	focused, err := xorg.QueryFocusedWindow(bridge)

	if err == nil && focused != nil {
		state.SetLastFocusedWindow(focused.ID)
	}

	screen := controllerScreen(config)

	if screen != nil {

		switch target.Position {

		case "right-of":
			bridge.WarpPointer(int(screen.Width)-2, int(screen.Height)/2)

		case "left-of":
			bridge.WarpPointer(1, int(screen.Height)/2)

		}

	}

	state.SetActive(target)

	xorg.UnfocusWindow(bridge)

	sendEnvelope(target, "focus", types.FocusEvent{
		Type: "focus",
		Edge: edge,
	})

	fmt.Printf("Activated remote machine for focus: %s (%s)\n", target.Hostname, target.Position)

	return true

}

func handoffToController(bridge *xorg.Bridge, state *types.GlobalState, config *types.Config, direction string) {

	state.ResetActive()

	screen := controllerScreen(config)

	if screen != nil {
		x := int(screen.Width) - 2
		if direction == "left" {
			x = 1
		}
		_, vy := state.GetCursor()
		bridge.WarpPointer(x, clamp(vy, 1, int(screen.Height)-2))
	}

	windows, err := xorg.QueryAllWindows(bridge)

	if err == nil && len(windows) > 0 {

		var target *types.Window

		if direction == "right" {
			target = helpers.FindLeftmostWindow(windows)
		} else {
			target = helpers.FindRightmostWindow(windows)
		}

		if target != nil {
			xorg.FocusWindow(bridge, target.ID)
		}

	}

	fmt.Printf("Handoff to controller: %s\n", direction)

}

func handleEventRequest(bridge *xorg.Bridge, state *types.GlobalState, config *types.Config, response http.ResponseWriter, request *http.Request) {

	body, err := io.ReadAll(request.Body)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	var event types.HandoffEvent

	if err := json.Unmarshal(body, &event); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "handoff" {

		active := state.GetActive()

		if active != nil && active.Hostname == event.Machine {
			handoffToController(bridge, state, config, event.Direction)
		}

	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("{\"status\": \"ok\"}"))

}
