//go:build integration

package xorg

import "testing"

func TestMoveWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-move", 0, 0, 300, 200)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	err = MoveWindow(bridge, id, 150, 75)

	if err != nil {
		t.Errorf("MoveWindow failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.X != 150 {
				t.Errorf("Expected X=%d, got %d", 150, w.X)
			}

			if w.Y != 75 {
				t.Errorf("Expected Y=%d, got %d", 75, w.Y)
			}

			return
		}
	}

	t.Errorf("Window not found after MoveWindow")

}
