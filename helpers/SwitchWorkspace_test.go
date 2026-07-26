package helpers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestSwitchWorkspace_NilBridge(t *testing.T) {

	state := types.NewGlobalState()

	err := SwitchWorkspace(nil, state, "FG")

	if err == nil {
		t.Error("Expected error with nil bridge")
	}

}

func TestSwitchWorkspace_NilState(t *testing.T) {

	err := SwitchWorkspace(nil, nil, "FG")

	if err == nil {
		t.Error("Expected error with nil state")
	}

}

func TestSwitchWorkspace_InvalidWorkspace(t *testing.T) {

	state := types.NewGlobalState()

	err := SwitchWorkspace(nil, state, "nonexistent")

	if err == nil {
		t.Error("Expected error for nonexistent workspace")
	}

}
