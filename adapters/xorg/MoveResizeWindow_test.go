//go:build integration

package xorg

import "testing"

func TestMoveResizeWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-move-resize", 0, 0, 300, 200)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	err = MoveResizeWindow(bridge, id, 100, 50, 500, 400)

	if err != nil {
		t.Errorf("MoveResizeWindow failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.X != 100 {
				t.Errorf("Expected X=%d, got %d", 100, w.X)
			}

			if w.Y != 50 {
				t.Errorf("Expected Y=%d, got %d", 50, w.Y)
			}

			if w.Width != 500 {
				t.Errorf("Expected Width=%d, got %d", 500, w.Width)
			}

			if w.Height != 400 {
				t.Errorf("Expected Height=%d, got %d", 400, w.Height)
			}

			return
		}
	}

	t.Errorf("Window not found after MoveResizeWindow")

}
