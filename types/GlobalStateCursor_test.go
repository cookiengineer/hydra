package types

import "testing"

func TestGlobalStateCursor(t *testing.T) {

	state := NewGlobalState()

	x, y := state.GetCursor()

	if x != 0 || y != 0 {
		t.Errorf("Expected initial cursor 0,0 got %d,%d", x, y)
	}

	state.SetCursor(3840+100, 540)

	x, y = state.GetCursor()

	if x != 3940 || y != 540 {
		t.Errorf("Expected cursor 3940,540 got %d,%d", x, y)
	}

}
