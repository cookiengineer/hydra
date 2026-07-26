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
	}

	for _, action := range expected {

		if !actions[action] {
			t.Errorf("Expected action %s in default bindings", action)
		}

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
