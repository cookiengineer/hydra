package types

import "testing"

func TestMachineParse(t *testing.T) {

	machine := &Machine{
		Hostname: " test-host ",
		IP:       "  192.168.1.1  ",
		Position: " left-of ",
		Screen: &Screen{
			Width:  1920,
			Height: 1080,
			Monitors: []Monitor{
				{Output: "HDMI-0", Connected: true, Resolution: "1920x1080", Width: 1920, Height: 1080},
			},
		},
	}

	err := machine.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if machine.Hostname != "test-host" {
		t.Errorf("Expected hostname test-host, got %s", machine.Hostname)
	}

	if machine.Position != "left-of" {
		t.Errorf("Expected position left-of, got %s", machine.Position)
	}

	if machine.Socket == nil {
		t.Error("Expected Socket channel to be initialized")
	}

}

func TestMachineParseInvalidPosition(t *testing.T) {

	machine := &Machine{
		Hostname: "test-host",
		IP:       "192.168.1.1",
		Position: "invalid-position",
		Screen: &Screen{
			Width:  1920,
			Height: 1080,
			Monitors: []Monitor{
				{Width: 1920, Height: 1080},
			},
		},
	}

	err := machine.Parse()

	if err == nil {
		t.Error("Expected error for invalid position")
	}

}

func TestMachineParseInvalidIP(t *testing.T) {

	machine := &Machine{
		Hostname: "test-host",
		IP:       "",
		Position: "center",
		Screen: &Screen{
			Width:  1920,
			Height: 1080,
			Monitors: []Monitor{
				{Width: 1920, Height: 1080},
			},
		},
	}

	err := machine.Parse()

	if err == nil {
		t.Error("Expected error for invalid IP")
	}

}

func TestMachineParseInvalidScreen(t *testing.T) {

	machine := &Machine{
		Hostname: "test-host",
		IP:       "192.168.1.1",
		Position: "center",
		Screen:   nil,
	}

	err := machine.Parse()

	if err == nil {
		t.Error("Expected error for nil screen")
	}

}

func TestMachineParseAllPositions(t *testing.T) {

	positions := []string{"left-of", "right-of", "above", "below", "center"}

	for _, position := range positions {

		machine := &Machine{
			Hostname: "test-host",
			IP:       "192.168.1.1",
			Position: position,
			Screen: &Screen{
				Width:  1920,
				Height: 1080,
				Monitors: []Monitor{
					{Width: 1920, Height: 1080},
				},
			},
		}

		err := machine.Parse()

		if err != nil {
			t.Errorf("Expected no error for position %s, got: %v", position, err)
		}

	}

}
