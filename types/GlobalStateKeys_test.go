package types

import "testing"

func TestGlobalStateTrackKeys(t *testing.T) {

	state := NewGlobalState()

	state.TrackKey(0xFFEB)
	state.TrackKey(0x61)
	state.TrackKey(0xFFEB)

	if !state.IsTrackedKey(0xFFEB) {
		t.Error("Expected Super_L to be tracked")
	}

	if !state.IsTrackedKey(0x61) {
		t.Error("Expected 'a' to be tracked")
	}

	if state.IsTrackedKey(0xFFE1) {
		t.Error("Did not expect Shift_L to be tracked")
	}

	if len(state.TrackedKeys()) != 2 {
		t.Errorf("Expected 2 tracked keys, got %d", len(state.TrackedKeys()))
	}

	state.UntrackKey(0xFFEB)

	if state.IsTrackedKey(0xFFEB) {
		t.Error("Expected Super_L to be untracked")
	}

	state.ClearTrackedKeys()

	if len(state.TrackedKeys()) != 0 {
		t.Errorf("Expected 0 tracked keys, got %d", len(state.TrackedKeys()))
	}

}

func TestGlobalStateTrackButtons(t *testing.T) {

	state := NewGlobalState()

	state.TrackButton(1)
	state.TrackButton(3)
	state.TrackButton(1)

	if !state.IsTrackedButton(1) {
		t.Error("Expected button 1 to be tracked")
	}

	if !state.IsTrackedButton(3) {
		t.Error("Expected button 3 to be tracked")
	}

	if state.IsTrackedButton(2) {
		t.Error("Did not expect button 2 to be tracked")
	}

	if len(state.TrackedButtons()) != 2 {
		t.Errorf("Expected 2 tracked buttons, got %d", len(state.TrackedButtons()))
	}

	state.UntrackButton(1)

	if state.IsTrackedButton(1) {
		t.Error("Expected button 1 to be untracked")
	}

	state.ClearTrackedButtons()

	if len(state.TrackedButtons()) != 0 {
		t.Errorf("Expected 0 tracked buttons, got %d", len(state.TrackedButtons()))
	}

}
