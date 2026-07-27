//go:build integration

package xorg

import "testing"

func TestSimulateKeyPress(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateKeyPress(bridge, 42)

	if err != nil {
		t.Errorf("SimulateKeyPress failed: %v", err)
	}

}
