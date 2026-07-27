//go:build integration

package xorg

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestTileWindow_Left(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-tile-left", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	monitor := &types.Monitor{
		Width:   1920,
		Height:  1080,
		OffsetX: 0,
		OffsetY: 0,
	}

	err = TileWindow(bridge, id, types.TileLeft, monitor)

	if err != nil {
		t.Errorf("TileWindow(Left) failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.X != 0 {
				t.Errorf("Expected X=%d, got %d", 0, w.X)
			}

			if w.Y != 0 {
				t.Errorf("Expected Y=%d, got %d", 0, w.Y)
			}

			if w.Width != 960 {
				t.Errorf("Expected Width=%d, got %d", 960, w.Width)
			}

			if w.Height != 1080 {
				t.Errorf("Expected Height=%d, got %d", 1080, w.Height)
			}

			return
		}
	}

	t.Errorf("Window not found after TileWindow(Left)")

}

func TestTileWindow_Right(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-tile-right", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	monitor := &types.Monitor{
		Width:   1920,
		Height:  1080,
		OffsetX: 0,
		OffsetY: 0,
	}

	err = TileWindow(bridge, id, types.TileRight, monitor)

	if err != nil {
		t.Errorf("TileWindow(Right) failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.X != 960 {
				t.Errorf("Expected X=%d, got %d", 960, w.X)
			}

			if w.Y != 0 {
				t.Errorf("Expected Y=%d, got %d", 0, w.Y)
			}

			if w.Width != 960 {
				t.Errorf("Expected Width=%d, got %d", 960, w.Width)
			}

			if w.Height != 1080 {
				t.Errorf("Expected Height=%d, got %d", 1080, w.Height)
			}

			return
		}
	}

	t.Errorf("Window not found after TileWindow(Right)")

}

func TestTileWindow_Top(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-tile-top", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	monitor := &types.Monitor{
		Width:   1920,
		Height:  1080,
		OffsetX: 0,
		OffsetY: 0,
	}

	err = TileWindow(bridge, id, types.TileTop, monitor)

	if err != nil {
		t.Errorf("TileWindow(Top) failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.X != 0 {
				t.Errorf("Expected X=%d, got %d", 0, w.X)
			}

			if w.Y != 0 {
				t.Errorf("Expected Y=%d, got %d", 0, w.Y)
			}

			if w.Width != 1920 {
				t.Errorf("Expected Width=%d, got %d", 1920, w.Width)
			}

			if w.Height != 540 {
				t.Errorf("Expected Height=%d, got %d", 540, w.Height)
			}

			return
		}
	}

	t.Errorf("Window not found after TileWindow(Top)")

}

func TestTileWindow_Bottom(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-tile-bottom", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	monitor := &types.Monitor{
		Width:   1920,
		Height:  1080,
		OffsetX: 0,
		OffsetY: 0,
	}

	err = TileWindow(bridge, id, types.TileBottom, monitor)

	if err != nil {
		t.Errorf("TileWindow(Bottom) failed: %v", err)
		return
	}

	windows, err := QueryAllWindows(bridge)

	if err != nil {
		t.Fatalf("QueryAllWindows failed: %v", err)
	}

	for _, w := range windows {
		if w.ID == id {
			if w.X != 0 {
				t.Errorf("Expected X=%d, got %d", 0, w.X)
			}

			if w.Y != 540 {
				t.Errorf("Expected Y=%d, got %d", 540, w.Y)
			}

			if w.Width != 1920 {
				t.Errorf("Expected Width=%d, got %d", 1920, w.Width)
			}

			if w.Height != 540 {
				t.Errorf("Expected Height=%d, got %d", 540, w.Height)
			}

			return
		}
	}

	t.Errorf("Window not found after TileWindow(Bottom)")

}

func TestTileWindow_Corners(t *testing.T) {

	tests := []struct {
		name     string
		position types.TilePosition
		expectedX int
		expectedY int
		expectedW int
		expectedH int
	}{
		{"TopLeft", types.TileTopLeft, 0, 0, 960, 540},
		{"TopRight", types.TileTopRight, 960, 0, 960, 540},
		{"BottomLeft", types.TileBottomLeft, 0, 540, 960, 540},
		{"BottomRight", types.TileBottomRight, 960, 540, 960, 540},
	}

	monitor := &types.Monitor{
		Width:   1920,
		Height:  1080,
		OffsetX: 0,
		OffsetY: 0,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bridge := createTestBridge(t)

			id, err := createTestWindow(bridge, "test-tile-"+tt.name, 0, 0, 400, 300)

			if err != nil {
				t.Fatalf("createTestWindow failed: %v", err)
			}

			defer destroyTestWindow(bridge, id)

			err = TileWindow(bridge, id, tt.position, monitor)

			if err != nil {
				t.Errorf("TileWindow(%s) failed: %v", tt.name, err)
				return
			}

			windows, err := QueryAllWindows(bridge)

			if err != nil {
				t.Fatalf("QueryAllWindows failed: %v", err)
			}

			found := false

			for _, w := range windows {
				if w.ID == id {
					found = true

					if w.X != tt.expectedX {
						t.Errorf("Expected X=%d, got %d", tt.expectedX, w.X)
					}

					if w.Y != tt.expectedY {
						t.Errorf("Expected Y=%d, got %d", tt.expectedY, w.Y)
					}

					if w.Width != tt.expectedW {
						t.Errorf("Expected Width=%d, got %d", tt.expectedW, w.Width)
					}

					if w.Height != tt.expectedH {
						t.Errorf("Expected Height=%d, got %d", tt.expectedH, w.Height)
					}

					break
				}
			}

			if !found {
				t.Errorf("Window not found after TileWindow(%s)", tt.name)
			}

		})
	}

}

func TestTileWindow_NilMonitor(t *testing.T) {

	bridge := createTestBridge(t)

	id, err := createTestWindow(bridge, "test-tile-nil-monitor", 0, 0, 400, 300)

	if err != nil {
		t.Fatalf("createTestWindow failed: %v", err)
	}

	defer destroyTestWindow(bridge, id)

	err = TileWindow(bridge, id, types.TileLeft, nil)

	if err == nil {
		t.Errorf("Expected error from TileWindow with nil monitor, got nil")
	}

}
