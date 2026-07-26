package helpers

import "fmt"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/types"

func MoveWindowToWorkspace(bridge *xorg.Bridge, state *types.GlobalState, targetName string) error {

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

	window, err := xorg.QueryFocusedWindow(bridge)

	if err != nil || window == nil {
		return fmt.Errorf("No focused window to move")
	}

	if current_ws != nil {

		current_windows := state.GetWorkspaceLayout(current_ws.Name)
		filtered := make([]types.Window, 0, len(current_windows))

		for _, win := range current_windows {

			if win.ID != window.ID {
				filtered = append(filtered, win)
			}

		}

		state.StoreWorkspaceLayout(current_ws.Name, filtered)

	}

	target_windows := state.GetWorkspaceLayout(target_ws.Name)
	target_windows = append(target_windows, *window)
	state.StoreWorkspaceLayout(target_ws.Name, target_windows)

	xorg.UnmapWindow(bridge, window.ID)

	fmt.Printf("Workspaces: Moved window %s to %s\n", window.Title, target_ws.Name)

	return nil

}
