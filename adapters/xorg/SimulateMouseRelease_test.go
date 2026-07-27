//go:build integration

package xorg

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestSimulateMouseRelease(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateMouseRelease(bridge, 500, 400, types.MouseButtonLeft)

	if err != nil {
		t.Errorf("SimulateMouseRelease failed: %v", err)
		return
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 500 || y != 400 {
		t.Errorf("Expected pointer at (500, 400), got (%d, %d)", x, y)
	}

}
