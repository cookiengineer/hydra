package types

import "encoding/json"
import "testing"

func TestNewEvent(t *testing.T) {

	payload := FocusEvent{Type: "focus", Direction: "left"}

	event, err := NewEvent("focus", payload)

	if err != nil {
		t.Fatalf("NewEvent failed: %v", err)
	}

	if event.Type != "focus" {
		t.Errorf("Expected type focus, got %s", event.Type)
	}

	data, err := json.Marshal(event)

	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed Event

	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if parsed.Type != "focus" {
		t.Errorf("Expected type focus, got %s", parsed.Type)
	}

	var decoded FocusEvent

	if err := parsed.Unmarshal(&decoded); err != nil {
		t.Fatalf("Event.Unmarshal failed: %v", err)
	}

	if decoded.Direction != "left" {
		t.Errorf("Expected direction left, got %s", decoded.Direction)
	}

}

func TestNewEventNilPayload(t *testing.T) {

	event, err := NewEvent("reset", nil)

	if err != nil {
		t.Fatalf("NewEvent failed: %v", err)
	}

	if len(event.Data) != 0 {
		t.Errorf("Expected empty data, got %s", string(event.Data))
	}

}

func TestKeyPressZeroValueDistinctFromUnset(t *testing.T) {

	press := KeyboardEvent{Type: KeyPress, Keycode: 42, Keysym: 0xFF53}

	data, err := json.Marshal(press)

	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed KeyboardEvent

	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if parsed.Type != KeyPress {
		t.Errorf("Expected KeyPress, got %v", parsed.Type)
	}

	envelope, err := NewEvent("keyboard", parsed)

	if err != nil {
		t.Fatalf("NewEvent failed: %v", err)
	}

	if envelope.Type != "keyboard" {
		t.Errorf("Expected envelope type keyboard, got %s", envelope.Type)
	}

}
