package types

import "encoding/json"
import "testing"

func TestResetEventJSON(t *testing.T) {

	event := ResetEvent{
		Type: "reset",
	}

	data, err := json.Marshal(event)

	if err != nil {
		t.Errorf("Failed to marshal ResetEvent: %s", err.Error())
	}

	var parsed ResetEvent

	err = json.Unmarshal(data, &parsed)

	if err != nil {
		t.Errorf("Failed to unmarshal ResetEvent: %s", err.Error())
	}

	if parsed.Type != "reset" {
		t.Errorf("Expected type reset, got %s", parsed.Type)
	}

}
