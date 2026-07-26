package types

import "testing"

func TestGlobalState(t *testing.T) {

	state := NewGlobalState()

	if state.GetActive() != nil {
		t.Error("Expected nil active initially")
	}

	machine := &Machine{Hostname: "test", Position: "left-of"}
	state.SetActive(machine)

	active := state.GetActive()

	if active == nil {
		t.Error("Expected active machine")
	}

	if active.Hostname != "test" {
		t.Errorf("Expected test, got %s", active.Hostname)
	}

	state.ResetActive()

	if state.GetActive() != nil {
		t.Error("Expected nil after reset")
	}

}

func TestVirtualScreen(t *testing.T) {

	vs := NewVirtualScreen()

	if vs == nil {
		t.Error("Expected non-nil virtual screen")
	}

	screen := &Screen{Width: 1920, Height: 1080}
	vs.AddMachine("test-host", screen)

	result := vs.GetMachine("test-host")

	if result == nil {
		t.Error("Expected to find test-host")
	}

	if result.Width != 1920 {
		t.Errorf("Expected width 1920, got %d", result.Width)
	}

	none := vs.GetMachine("nonexistent")

	if none != nil {
		t.Error("Expected nil for nonexistent host")
	}

}
