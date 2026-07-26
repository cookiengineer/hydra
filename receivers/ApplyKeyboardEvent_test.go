package receivers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestApplyKeyboardEvent_NoopWithNilBridge(t *testing.T) {

	event := &types.KeyboardEvent{
		Type:    types.KeyPress,
		Keycode: 42,
	}

	// Should not panic with nil bridge
	ApplyKeyboardEvent(nil, event)

	if event.Type != types.KeyPress {
		t.Errorf("Expected KeyPress type, got %v", event.Type)
	}

	if event.Keycode != 42 {
		t.Errorf("Expected keycode 42, got %d", event.Keycode)
	}

}

func TestApplyKeyboardEvent_AllTypes(t *testing.T) {

	tests := []struct {
		name    string
		event   types.KeyboardEvent
	}{
		{"KeyPress", types.KeyboardEvent{Type: types.KeyPress, Keycode: 10}},
		{"KeyRelease", types.KeyboardEvent{Type: types.KeyRelease, Keycode: 20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyKeyboardEvent(nil, &tt.event)
		})
	}

}
