//go:build integration

package xorg

import "testing"

func TestQueryAllWindows_SeesNewWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-sees-new", 50, 50, 300, 200)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	found := false

	for _, w := range windows {
		if w.ID == id {
			found = true

			if w.X != 50 {
				t.Errorf("Expected X=%d, got %d", 50, w.X)
			}

			if w.Y != 50 {
				t.Errorf("Expected Y=%d, got %d", 50, w.Y)
			}

			if w.Width != 300 {
				t.Errorf("Expected Width=%d, got %d", 300, w.Width)
			}

			if w.Height != 200 {
				t.Errorf("Expected Height=%d, got %d", 200, w.Height)
			}

			if w.Title != "test-sees-new" {
				t.Errorf("Expected Title=%q, got %q", "test-sees-new", w.Title)
			}

			break
		}
	}

	if !found {
		t.Errorf("Created window ID=%d not found in QueryAllWindows", id)
	}

}
