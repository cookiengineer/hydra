//go:build integration

package xorg

import "testing"

func TestSimulateMouseScroll(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateMouseScroll(bridge, 600, 500, 0, 3)

	if err != nil {
		t.Errorf("SimulateMouseScroll failed: %v", err)
		return
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 600 || y != 500 {
		t.Errorf("Expected pointer at (600, 500), got (%d, %d)", x, y)
	}

}

func TestSimulateMouseScroll_NegativeDY(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateMouseScroll(bridge, 700, 600, 0, -2)

	if err != nil {
		t.Errorf("SimulateMouseScroll failed: %v", err)
		return
	}

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 700 || y != 600 {
		t.Errorf("Expected pointer at (700, 600), got (%d, %d)", x, y)
	}

}
