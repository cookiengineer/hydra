package types

import "encoding/json"
import "testing"

func TestInitEventJSON(t *testing.T) {

	ws := Workspace{Name: "FG", Index: 0}
	vs := &VirtualScreen{Width: 1920, Height: 1080}

	event := InitEvent{
		Type:            "init",
		VirtualScreen:   vs,
		Workspaces:      []Workspace{ws},
		ActiveWorkspace: "FG",
	}

	data, err := json.Marshal(event)

	if err != nil {
		t.Errorf("Failed to marshal InitEvent: %s", err.Error())
	}

	var parsed InitEvent

	err = json.Unmarshal(data, &parsed)

	if err != nil {
		t.Errorf("Failed to unmarshal InitEvent: %s", err.Error())
	}

	if parsed.Type != "init" {
		t.Errorf("Expected type init, got %s", parsed.Type)
	}

	if parsed.ActiveWorkspace != "FG" {
		t.Errorf("Expected active_workspace FG, got %s", parsed.ActiveWorkspace)
	}

	if len(parsed.Workspaces) != 1 {
		t.Errorf("Expected 1 workspace, got %d", len(parsed.Workspaces))
	}

	if parsed.Workspaces[0].Name != "FG" {
		t.Errorf("Expected FG, got %s", parsed.Workspaces[0].Name)
	}

	if parsed.VirtualScreen.Width != 1920 {
		t.Errorf("Expected width 1920, got %d", parsed.VirtualScreen.Width)
	}

}
