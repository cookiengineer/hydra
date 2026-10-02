package helpers

import "testing"
import "github.com/cookiengineer/hydra/types"

func edgeWindows() []types.Window {

	return []types.Window{
		{ID: 1, Title: "a", X: 500, Y: 0, Width: 200, Height: 200},
		{ID: 2, Title: "b", X: 100, Y: 300, Width: 200, Height: 200},
		{ID: 3, Title: "c", X: 900, Y: 300, Width: 200, Height: 200},
		{ID: 4, Title: "d", X: 100, Y: 100, Width: 200, Height: 200},
	}

}

func TestFindLeftmostWindow(t *testing.T) {

	windows := edgeWindows()
	result := FindLeftmostWindow(windows)

	if result == nil || result.ID != 4 {
		t.Errorf("Expected window 4 (top-most tie break), got %+v", result)
	}

}

func TestFindRightmostWindow(t *testing.T) {

	windows := edgeWindows()
	result := FindRightmostWindow(windows)

	if result == nil || result.ID != 3 {
		t.Errorf("Expected window 3, got %+v", result)
	}

}

func TestFindTopmostWindow(t *testing.T) {

	windows := edgeWindows()
	result := FindTopmostWindow(windows)

	if result == nil || result.ID != 1 {
		t.Errorf("Expected window 1, got %+v", result)
	}

}

func TestFindBottommostWindow(t *testing.T) {

	windows := edgeWindows()
	result := FindBottommostWindow(windows)

	if result == nil || result.ID != 2 {
		t.Errorf("Expected window 2, got %+v", result)
	}

}

func TestFindWindowEdgeEmpty(t *testing.T) {

	if FindLeftmostWindow(nil) != nil {
		t.Error("Expected nil for empty windows")
	}

	if FindRightmostWindow(nil) != nil {
		t.Error("Expected nil for empty windows")
	}

}
