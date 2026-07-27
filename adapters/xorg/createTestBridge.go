//go:build integration

package xorg

import "os"
import "testing"

func createTestBridge(t *testing.T) *Bridge {

	display := os.Getenv("DISPLAY")

	if display == "" {
		t.Skip("DISPLAY not set; start Xvfb: Xvfb :99 -screen 0 1920x1080x24 &")
	}

	bridge, err := NewBridge(display)

	if err != nil {
		t.Fatalf("NewBridge(%q) failed: %v", display, err)
	}

	t.Cleanup(func() {
		bridge.Destroy()
	})

	return bridge

}

