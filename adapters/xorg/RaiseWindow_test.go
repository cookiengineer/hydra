//go:build integration

package xorg

import "testing"

func TestRaiseWindow(t *testing.T) {

	bridge := createTestBridge(t)

	id1, err := createTestWindow(bridge, "test-raise-1", 10, 10, 200, 100)

	if err != nil {
		t.Fatalf("createTestWindow for id1 failed: %v", err)
	}

	defer destroyTestWindow(bridge, id1)

	id2, err := createTestWindow(bridge, "test-raise-2", 30, 30, 200, 100)

	if err != nil {
		t.Fatalf("createTestWindow for id2 failed: %v", err)
	}

	defer destroyTestWindow(bridge, id2)

	err = RaiseWindow(bridge, id1)

	if err != nil {
		t.Errorf("RaiseWindow failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	var id1_index int = -1
	var id2_index int = -1

	for i, w := range windows {
		if w.ID == id1 {
			id1_index = i
		}

		if w.ID == id2 {
			id2_index = i
		}
	}

	if id1_index == -1 || id2_index == -1 {
		t.Errorf("One or both windows not found in QueryAllWindows")
		return
	}

	if id1_index <= id2_index {
		t.Errorf("Raised window ID=%d should be above window ID=%d in stacking order", id1, id2)
	}

}
