//go:build integration

package xorg

import "testing"

func TestResizeWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-resize", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	err = ResizeWindow(bridge, id, 640, 480)

	if err != nil {
		t.Errorf("ResizeWindow failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.Width != 640 {
				t.Errorf("Expected Width=%d, got %d", 640, w.Width)
			}

			if w.Height != 480 {
				t.Errorf("Expected Height=%d, got %d", 480, w.Height)
			}

			return
		}
	}

	t.Errorf("Window not found after ResizeWindow")

}
