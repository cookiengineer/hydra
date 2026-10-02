package actions

import "testing"

func TestClamp(t *testing.T) {

	if clamp(5, 0, 10) != 5 {
		t.Error("Expected 5 within range")
	}

	if clamp(-3, 0, 10) != 0 {
		t.Error("Expected clamp to minimum")
	}

	if clamp(42, 0, 10) != 10 {
		t.Error("Expected clamp to maximum")
	}

}
