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

func TestButtonMask(t *testing.T) {

	if buttonMask(1) != 256 {
		t.Errorf("Expected Button1Mask 256, got %d", buttonMask(1))
	}

	if buttonMask(2) != 512 {
		t.Errorf("Expected Button2Mask 512, got %d", buttonMask(2))
	}

	if buttonMask(3) != 1024 {
		t.Errorf("Expected Button3Mask 1024, got %d", buttonMask(3))
	}

	if buttonMask(0) != 0 {
		t.Errorf("Expected 0 for invalid button, got %d", buttonMask(0))
	}

}
