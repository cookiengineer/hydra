//go:build integration

package xorg

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestSimulateMousePress(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateMousePress(bridge, 400, 300, types.MouseButtonLeft)

	if err != nil {
		t.Errorf("SimulateMousePress failed: %v", err)
		return
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 400 || y != 300 {
		t.Errorf("Expected pointer at (400, 300), got (%d, %d)", x, y)
	}

}
