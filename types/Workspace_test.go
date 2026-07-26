package types

import "testing"

func TestWorkspace(t *testing.T) {

	ws := Workspace{
		Name:    "FG",
		Index:   0,
		Windows: []Window{},
	}

	if ws.Name != "FG" {
		t.Errorf("Expected FG, got %s", ws.Name)
	}

	if ws.Index != 0 {
		t.Errorf("Expected 0, got %d", ws.Index)
	}

	if ws.Windows == nil {
		t.Error("Expected non-nil Windows slice")
	}

}
