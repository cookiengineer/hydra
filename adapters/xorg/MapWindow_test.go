//go:build integration

package xorg

import "testing"

func TestMapUnmapWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-map-unmap", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	err = UnmapWindow(bridge, id)

	if err != nil {
		t.Errorf("UnmapWindow failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			t.Errorf("Window should not appear in QueryAllWindows after UnmapWindow")
			return
		}
	}

	err = MapWindow(bridge, id)

	if err != nil {
		t.Errorf("MapWindow failed: %v", err)
		return
	}

	windows, err = QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	found := false

	for _, w := range windows {
		if w.ID == id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Window should appear in QueryAllWindows after MapWindow")
	}

}
