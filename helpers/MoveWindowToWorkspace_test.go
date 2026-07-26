package helpers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestMoveWindowToWorkspace_NilBridge(t *testing.T) {

	state := types.NewGlobalState()

	err := MoveWindowToWorkspace(nil, state, "FG")

	if err == nil {
		t.Error("Expected error with nil bridge")
	}

}

func TestMoveWindowToWorkspace_NilState(t *testing.T) {

	err := MoveWindowToWorkspace(nil, nil, "FG")

	if err == nil {
		t.Error("Expected error with nil state")
	}

}

func TestMoveWindowToWorkspace_InvalidWorkspace(t *testing.T) {

	state := types.NewGlobalState()

	err := MoveWindowToWorkspace(nil, state, "nonexistent")

	if err == nil {
		t.Error("Expected error for nonexistent workspace")
	}

}
