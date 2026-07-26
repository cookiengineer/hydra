package helpers

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func SwitchWorkspace(bridge *xorg.Bridge, state *types.GlobalState, targetName string) error {

	if bridge == nil {
		return fmt.Errorf("Bridge is nil")
	}

	if state == nil {
		return fmt.Errorf("State is nil")
	}

	current_ws := state.GetWorkspaceByName(state.GetActiveWorkspace())
	target_ws := state.GetWorkspaceByName(targetName)

	if target_ws == nil {
		return fmt.Errorf("Workspace %s not found", targetName)
	}

	windows, err := xorg.QueryAllWindows(bridge)

	if err == nil && current_ws != nil {

		state.StoreWorkspaceLayout(current_ws.Name, windows)

		for _, win := range windows {
			xorg.UnmapWindow(bridge, win.ID)
		}

	}

	target_windows := state.GetWorkspaceLayout(target_ws.Name)

	for _, win := range target_windows {
		xorg.MoveResizeWindow(bridge, win.ID, win.X, win.Y, win.Width, win.Height)
		xorg.MapWindow(bridge, win.ID)
	}

	if len(target_windows) > 0 {
		xorg.FocusWindow(bridge, target_windows[0].ID)
	}

	state.SetActiveWorkspace(target_ws.Name)

	fmt.Printf("Workspaces: Switched to %s\n", target_ws.Name)

	return nil

}
