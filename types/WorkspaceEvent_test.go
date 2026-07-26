package types

import "encoding/json"
import "testing"

func TestWorkspaceEventJSON(t *testing.T) {

	event := WorkspaceEvent{
		Type:  "workspace",
		Name:  "FG",
		Index: 0,
	}

	data, err := json.Marshal(event)

	if err != nil {
		t.Errorf("Failed to marshal WorkspaceEvent: %s", err.Error())
	}

	var parsed WorkspaceEvent

	err = json.Unmarshal(data, &parsed)

	if err != nil {
		t.Errorf("Failed to unmarshal WorkspaceEvent: %s", err.Error())
	}

	if parsed.Type != "workspace" {
		t.Errorf("Expected type workspace, got %s", parsed.Type)
	}

	if parsed.Name != "FG" {
		t.Errorf("Expected FG, got %s", parsed.Name)
	}

	if parsed.Index != 0 {
		t.Errorf("Expected index 0, got %d", parsed.Index)
	}

}
