package types

import "encoding/json"
import "testing"

func TestFocusEventJSON(t *testing.T) {

	tests := []struct {
		direction string
	}{
		{"left"},
		{"right"},
		{"up"},
		{"down"},
	}

	for _, test := range tests {

		event := FocusEvent{
			Type:      "focus",
			Direction: test.direction,
		}

		data, err := json.Marshal(event)

		if err != nil {
			t.Errorf("Failed to marshal FocusEvent: %s", err.Error())
		}

		var parsed FocusEvent

		err = json.Unmarshal(data, &parsed)

		if err != nil {
			t.Errorf("Failed to unmarshal FocusEvent: %s", err.Error())
		}

		if parsed.Type != "focus" {
			t.Errorf("Expected type focus, got %s", parsed.Type)
		}

		if parsed.Direction != test.direction {
			t.Errorf("Expected direction %s, got %s", test.direction, parsed.Direction)
		}

	}

}
