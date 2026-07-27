//go:build integration

package xorg

import "testing"

func TestFocusWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id1, err := createTestWindow(bridge, "test-focus-1", 10, 10, 200, 100)

	if err != nil {
		t.Fatalf("createTestWindow for id1 failed: %v", err)
	}

	defer destroyTestWindow(bridge, id1)

	id2, err := createTestWindow(bridge, "test-focus-2", 300, 10, 200, 100)

	if err != nil {
		t.Fatalf("createTestWindow for id2 failed: %v", err)
	}

	defer destroyTestWindow(bridge, id2)

	err = FocusWindow(bridge, id1)

	if err != nil {
		t.Errorf("FocusWindow failed: %v", err)
		return
	}

	focused, err := QueryFocusedWindow(bridge)

	if err != nil {
		t.Errorf("QueryFocusedWindow failed: %v", err)
		return
	}

	if focused.ID != id1 {
		t.Errorf("Expected focused window ID=%d, got ID=%d", id1, focused.ID)
	}

}
