package types

import "testing"

func TestKeyBindingMatches(t *testing.T) {

	binding := KeyBinding{
		Modifiers: SuperKeyMask,
		Keycode:   XK_Escape,
		Action:    ActionResetToController,
	}

	if !binding.Matches(SuperKeyMask, XK_Escape) {
		t.Error("Expected match")
	}

	if binding.Matches(SuperKeyMask, XK_Left) {
		t.Error("Expected no match for different keycode")
	}

	if binding.Matches(0, XK_Escape) {
		t.Error("Expected no match for different modifiers")
	}

}

func TestGetDefaultKeyBindings(t *testing.T) {

	bindings := GetDefaultKeyBindings()

	if len(bindings) == 0 {
		t.Error("Expected non-empty bindings list")
	}

	actions := make(map[string]bool)

	for _, b := range bindings {
		actions[b.Action] = true
	}

	expected := []string{
		ActionResetToController,
		ActionFocusLeft,
		ActionFocusRight,
		ActionFocusUp,
		ActionFocusDown,
		ActionTileLeft,
		ActionTileRight,
		ActionTileTop,
		ActionTileBottom,
		ActionSwitchWorkspace,
		ActionMoveToWorkspace,
	}

	for _, action := range expected {

		if !actions[action] {
			t.Errorf("Expected action %s in default bindings", action)
		}

	}

}

func TestDefaultBindingCount(t *testing.T) {

	bindings := GetDefaultKeyBindings()

	expected := 9 + 14 + 14 // 9 base + 14 switch + 14 move

	if len(bindings) != expected {
		t.Errorf("Expected %d default bindings, got %d", expected, len(bindings))
	}

}

func TestWorkspaceBindingData(t *testing.T) {

	bindings := GetDefaultKeyBindings()

	workspace_switch_count := 0
	workspace_move_count := 0

	for _, b := range bindings {

		if b.Action == ActionSwitchWorkspace {
			workspace_switch_count++
			if b.Data > 13 {
				t.Errorf("Expected Data <= 13, got %d", b.Data)
			}
			if b.Modifiers != SuperKeyMask {
				t.Errorf("Expected SuperKeyMask only, got %d", b.Modifiers)
			}
		}

		if b.Action == ActionMoveToWorkspace {
			workspace_move_count++
			if b.Data > 13 {
				t.Errorf("Expected Data <= 13, got %d", b.Data)
			}
			if b.Modifiers != SuperKeyMask|ModShift {
				t.Errorf("Expected SuperKeyMask|ModShift, got %d", b.Modifiers)
			}
		}

	}

	if workspace_switch_count != 14 {
		t.Errorf("Expected 14 workspace switch bindings, got %d", workspace_switch_count)
	}

	if workspace_move_count != 14 {
		t.Errorf("Expected 14 workspace move bindings, got %d", workspace_move_count)
	}

}

func TestKeyBindingData(t *testing.T) {

	binding := KeyBinding{
		Modifiers: SuperKeyMask,
		Keycode:   XK_grave,
		Action:    ActionSwitchWorkspace,
		Data:      0,
	}

	if binding.Data != 0 {
		t.Errorf("Expected Data 0, got %d", binding.Data)
	}

	if !binding.Matches(SuperKeyMask, XK_grave) {
		t.Error("Expected match for workspace binding")
	}

}

func TestTilePositionString(t *testing.T) {

	tests := []struct {
		position TilePosition
		expected string
	}{
		{TileLeft, "left"},
		{TileRight, "right"},
		{TileTop, "top"},
		{TileBottom, "bottom"},
		{TileTopLeft, "top-left"},
		{TileTopRight, "top-right"},
		{TileBottomLeft, "bottom-left"},
		{TileBottomRight, "bottom-right"},
	}

	for _, test := range tests {

		if test.position.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.position.String())
		}

	}

}

func TestTilePositionFromString(t *testing.T) {

	tests := []struct {
		input    string
		expected TilePosition
	}{
		{"left", TileLeft},
		{"right", TileRight},
		{"top", TileTop},
		{"bottom", TileBottom},
		{"top-left", TileTopLeft},
		{"top-right", TileTopRight},
		{"bottom-left", TileBottomLeft},
		{"bottom-right", TileBottomRight},
		{"invalid", -1},
	}

	for _, test := range tests {

		result := TilePositionFromString(test.input)

		if result != test.expected {
			t.Errorf("Expected %d, got %d for input %s", test.expected, result, test.input)
		}

	}

}
