//go:build integration

package xorg

import "testing"

func TestUnfocusWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-unfocus", 10, 10, 200, 100)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	err = FocusWindow(bridge, id)

	if err != nil {
		t.Errorf("FocusWindow failed: %v", err)
		return
	}

	focused, err := QueryFocusedWindow(bridge)

	if err != nil {
		t.Errorf("QueryFocusedWindow after FocusWindow failed: %v", err)
		return
	}

	if focused.ID != id {
		t.Errorf("Expected focused window ID=%d, got ID=%d", id, focused.ID)
		return
	}

	err = UnfocusWindow(bridge)

	if err != nil {
		t.Errorf("UnfocusWindow failed: %v", err)
		return
	}

	_, err = QueryFocusedWindow(bridge)

	if err == nil {
		t.Errorf("Expected error from QueryFocusedWindow after UnfocusWindow, got nil")
	}

}
