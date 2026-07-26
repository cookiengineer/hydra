package types

import "testing"

func TestNewConfig(t *testing.T) {

	machine := Machine{
		Hostname: "test-host",
		IP:       "192.168.1.1",
		Position: "center",
		Screen: &Screen{
			Width:  1920,
			Height: 1080,
			Monitors: []Monitor{
				{Output: "HDMI-0", Connected: true, Resolution: "1920x1080", Width: 1920, Height: 1080},
			},
		},
	}

	config := NewConfig(machine)

	if config.Controller != "test-host" {
		t.Errorf("Expected controller test-host, got %s", config.Controller)
	}

	if len(config.Machines) != 1 {
		t.Errorf("Expected 1 machine, got %d", len(config.Machines))
	}

	if config.Machines["test-host"] == nil {
		t.Error("Expected controller machine to be in Machines map")
	}

	if len(config.Workspaces) != 14 {
		t.Errorf("Expected 14 workspaces, got %d", len(config.Workspaces))
	}

	if len(config.KeyBindings) == 0 {
		t.Error("Expected non-empty KeyBindings")
	}

}

func TestGetDefaultWorkspaces(t *testing.T) {

	workspaces := GetDefaultWorkspaces()

	if len(workspaces) != 14 {
		t.Errorf("Expected 14 workspaces, got %d", len(workspaces))
	}

	if workspaces[0].Name != "FG" || workspaces[0].Index != 0 {
		t.Errorf("Expected workspace FG at index 0, got %s at %d", workspaces[0].Name, workspaces[0].Index)
	}

	if workspaces[13].Name != "BG" || workspaces[13].Index != 13 {
		t.Errorf("Expected workspace BG at index 13, got %s at %d", workspaces[13].Name, workspaces[13].Index)
	}

	if workspaces[1].Name != "1" {
		t.Errorf("Expected 1, got %s", workspaces[1].Name)
	}

	if workspaces[10].Name != "10" {
		t.Errorf("Expected 10, got %s", workspaces[10].Name)
	}

}

func TestGetMachine(t *testing.T) {

	config := NewConfig(Machine{
		Hostname: "controller",
		IP:       "10.0.0.1",
		Position: "center",
		Screen:   &Screen{Width: 1920, Height: 1080, Monitors: []Monitor{{Width: 1920, Height: 1080}}},
	})

	found := config.GetMachine("controller")

	if found == nil {
		t.Error("Expected to find controller machine")
	}

	not_found := config.GetMachine("nonexistent")

	if not_found != nil {
		t.Error("Expected nil for nonexistent machine")
	}

}

func TestQueryMachine(t *testing.T) {

	config := NewConfig(Machine{
		Hostname: "controller",
		IP:       "10.0.0.1",
		Position: "center",
		Screen:   &Screen{Width: 1920, Height: 1080, Monitors: []Monitor{{Width: 1920, Height: 1080}}},
	})

	config.UpdateMachine(Machine{
		Hostname: "left1",
		IP:       "10.0.0.2",
		Position: "left-of",
		Screen:   &Screen{Width: 1280, Height: 720, Monitors: []Monitor{{Width: 1280, Height: 720}}},
	})

	left := config.QueryMachine("left-of")

	if left == nil {
		t.Error("Expected to find left-of machine")
	}

	if left.Hostname != "left1" {
		t.Errorf("Expected left1, got %s", left.Hostname)
	}

	none := config.QueryMachine("right-of")

	if none != nil {
		t.Error("Expected nil for right-of position")
	}

}

func TestUpdateMachine(t *testing.T) {

	config := NewConfig(Machine{
		Hostname: "controller",
		IP:       "10.0.0.1",
		Position: "center",
		Screen:   &Screen{Width: 1920, Height: 1080, Monitors: []Monitor{{Width: 1920, Height: 1080}}},
	})

	ok := config.UpdateMachine(Machine{
		Hostname: "new-machine",
		IP:       "10.0.0.3",
		Position: "right-of",
		Screen:   &Screen{Width: 1280, Height: 720, Monitors: []Monitor{{Width: 1280, Height: 720}}},
	})

	if ok == false {
		t.Error("Expected UpdateMachine to return true")
	}

	if config.GetMachine("new-machine") == nil {
		t.Error("Expected new-machine to be in config")
	}

}

func TestRemoveMachine(t *testing.T) {

	config := NewConfig(Machine{
		Hostname: "controller",
		IP:       "10.0.0.1",
		Position: "center",
		Screen:   &Screen{Width: 1920, Height: 1080, Monitors: []Monitor{{Width: 1920, Height: 1080}}},
	})

	config.UpdateMachine(Machine{
		Hostname: "to-remove",
		IP:       "10.0.0.4",
		Position: "left-of",
		Screen:   &Screen{Width: 1280, Height: 720, Monitors: []Monitor{{Width: 1280, Height: 720}}},
	})

	ok := config.RemoveMachine(Machine{Hostname: "to-remove"})

	if !ok {
		t.Error("Expected RemoveMachine to return true")
	}

	if config.GetMachine("to-remove") != nil {
		t.Error("Expected to-remove to be gone")
	}

}

func TestSetThis(t *testing.T) {

	config := NewConfig(Machine{
		Hostname: "controller",
		IP:       "10.0.0.1",
		Position: "center",
		Screen:   &Screen{Width: 1920, Height: 1080, Monitors: []Monitor{{Width: 1920, Height: 1080}}},
	})

	ok1 := config.SetThis("controller")

	if ok1 == false {
		t.Error("Expected SetThis to return true")
	}

	if config.This != "controller" {
		t.Errorf("Expected This to be controller, got %s", config.This)
	}

	ok2 := config.SetThis("nonexistent")

	if ok2 == true {
		t.Error("Expected SetThis to return false for nonexistent machine")
	}

}
