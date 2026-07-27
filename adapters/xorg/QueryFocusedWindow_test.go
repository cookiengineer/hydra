//go:build integration

package xorg

import "testing"

func TestQueryFocusedWindow_NoFocus(t *testing.T) {

	bridge := createTestBridge(t)

	_, err := QueryFocusedWindow(bridge)

	if err == nil {
		t.Errorf("Expected error from QueryFocusedWindow when no window is focused, got nil")
	}

}
