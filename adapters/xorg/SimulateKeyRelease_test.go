//go:build integration

package xorg

import "testing"

func TestSimulateKeyRelease(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateKeyRelease(bridge, 42)

	if err != nil {
		t.Errorf("SimulateKeyRelease failed: %v", err)
	}

}
