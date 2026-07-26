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

func TestWorkspaceState(t *testing.T) {

	state := NewGlobalState()

	if state.GetActiveWorkspace() != "FG" {
		t.Errorf("Expected active workspace FG, got %s", state.GetActiveWorkspace())
	}

	if state.Workspaces == nil {
		t.Error("Expected non-nil Workspaces map")
	}

	if len(state.Workspaces) != 14 {
		t.Errorf("Expected 14 workspaces, got %d", len(state.Workspaces))
	}

	fg := state.GetWorkspace("FG")

	if fg == nil {
		t.Error("Expected FG workspace to exist")
	}

	if fg.Name != "FG" {
		t.Errorf("Expected FG, got %s", fg.Name)
	}

	if fg.Index != 0 {
		t.Errorf("Expected index 0, got %d", fg.Index)
	}

}

func TestSetActiveWorkspace(t *testing.T) {

	state := NewGlobalState()
	state.SetActiveWorkspace("5")

	if state.GetActiveWorkspace() != "5" {
		t.Errorf("Expected active workspace 5, got %s", state.GetActiveWorkspace())
	}

}

func TestStoreWorkspaceLayout(t *testing.T) {

	state := NewGlobalState()

	windows := []Window{
		{ID: 1, Title: "Terminal", X: 0, Y: 0, Width: 960, Height: 1080},
		{ID: 2, Title: "Browser", X: 960, Y: 0, Width: 960, Height: 1080},
	}

	state.StoreWorkspaceLayout("FG", windows)

	stored := state.GetWorkspaceLayout("FG")

	if len(stored) != 2 {
		t.Errorf("Expected 2 windows, got %d", len(stored))
	}

	if stored[0].ID != 1 || stored[1].ID != 2 {
		t.Error("Expected window IDs 1 and 2")
	}

}

func TestGetWorkspaceLayoutNonExistent(t *testing.T) {

	state := NewGlobalState()
	layout := state.GetWorkspaceLayout("nonexistent")

	if layout != nil {
		t.Error("Expected nil for nonexistent workspace")
	}

}

func TestGetWorkspaceByName(t *testing.T) {

	state := NewGlobalState()

	bg := state.GetWorkspaceByName("BG")

	if bg == nil {
		t.Error("Expected BG workspace to exist")
	}

	if bg.Index != 13 {
		t.Errorf("Expected index 13, got %d", bg.Index)
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
