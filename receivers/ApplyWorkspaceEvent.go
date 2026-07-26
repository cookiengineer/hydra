package receivers

import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/helpers"
import "github.com/cookiengineer/hydra/types"

func ApplyWorkspaceEvent(bridge *xorg.Bridge, state *types.GlobalState, event *types.WorkspaceEvent) {

	if bridge != nil && state != nil {
		helpers.SwitchWorkspace(bridge, state, event.Name)
	}

}
