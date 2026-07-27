//go:build integration

package xorg

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestSimulateMouseMove(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateMouseMove(bridge, 200, 150, 37, 37)

	if err != nil {
		t.Errorf("SimulateMouseMove failed: %v", err)
		return
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 237 || y != 187 {
		t.Errorf("Expected pointer at (237, 187), got (%d, %d)", x, y)
	}

}

func TestSimulateMouseMove_WithScreen(t *testing.T) {

	bridge := createTestBridge(t)
	bridge.Screen = &types.Screen{
		Width:   800,
		Height:  600,
		OffsetX: 100,
		OffsetY: 50,
	}

	err := SimulateMouseMove(bridge, 300, 200, 50, 30)

	if err != nil {
		t.Errorf("SimulateMouseMove failed: %v", err)
		return
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	expected_x := int(300 + 50 - 100)
	expected_y := int(200 + 30 - 50)

	if x != expected_x || y != expected_y {
		t.Errorf("Expected pointer at (%d, %d), got (%d, %d)", expected_x, expected_y, x, y)
	}

}
