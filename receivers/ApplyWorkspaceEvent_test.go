package receivers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestApplyWorkspaceEvent_NoopWithNilBridge(t *testing.T) {

	state := types.NewGlobalState()
	event := &types.WorkspaceEvent{
		Type:  "workspace",
		Name:  "FG",
		Index: 0,
	}

	ApplyWorkspaceEvent(nil, state, event)

}

func TestApplyWorkspaceEvent_NoopWithNilState(t *testing.T) {

	event := &types.WorkspaceEvent{
		Type:  "workspace",
		Name:  "1",
		Index: 1,
	}

	ApplyWorkspaceEvent(nil, nil, event)

}

func TestApplyWorkspaceEvent_AllWorkspaces(t *testing.T) {

	state := types.NewGlobalState()

	names := []string{"FG", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "BG"}

	for i, name := range names {

		event := &types.WorkspaceEvent{
			Type:  "workspace",
			Name:  name,
			Index: uint32(i),
		}

		ApplyWorkspaceEvent(nil, state, event)

	}

}
