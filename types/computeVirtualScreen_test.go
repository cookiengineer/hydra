package types

import "testing"

func makeMachine(hostname string, position string, width uint, height uint) *Machine {

	return &Machine{
		Hostname: hostname,
		IP:       "0.0.0.0",
		Position: position,
		Screen: &Screen{
			Width:    width,
			Height:   height,
			Monitors: []Monitor{
				{
					Output:     "HDMI-A-0",
					Connected:  true,
					Resolution: "1920x1080",
					Width:      int(width),
					Height:     int(height),
					OffsetX:    0,
					OffsetY:    0,
				},
			},
		},
	}

}

func TestComputeVirtualScreen_ControllerOnly(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)

	vs := computeVirtualScreen("controller", machines)

	if vs.Width != 1920 {
		t.Errorf("Expected virtual width 1920, got %d", vs.Width)
	}

	if vs.Height != 1080 {
		t.Errorf("Expected virtual height 1080, got %d", vs.Height)
	}

	ctrl_screen := vs.GetMachine("controller")

	if ctrl_screen == nil {
		t.Fatal("Expected controller in virtual screen")
	}

	if ctrl_screen.OffsetX != 0 {
		t.Errorf("Expected controller OffsetX 0, got %d", ctrl_screen.OffsetX)
	}

	if ctrl_screen.OffsetY != 0 {
		t.Errorf("Expected controller OffsetY 0, got %d", ctrl_screen.OffsetY)
	}

}

func TestComputeVirtualScreen_LeftOfController(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["left1"]      = makeMachine("left1", "left-of", 1280, 720)

	vs := computeVirtualScreen("controller", machines)

	expected_width := uint(1280 + 1920)

	if vs.Width != expected_width {
		t.Errorf("Expected virtual width %d, got %d", expected_width, vs.Width)
	}

	ctrl := vs.GetMachine("controller")
	left := vs.GetMachine("left1")

	if left.OffsetX != 0 {
		t.Errorf("Expected left1 OffsetX 0, got %d", left.OffsetX)
	}

	if ctrl.OffsetX != 1280 {
		t.Errorf("Expected controller OffsetX 1280, got %d", ctrl.OffsetX)
	}

}

func TestComputeVirtualScreen_RightOfController(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["right1"]     = makeMachine("right1", "right-of", 1280, 720)

	vs := computeVirtualScreen("controller", machines)

	expected_width := uint(1920 + 1280)

	if vs.Width != expected_width {
		t.Errorf("Expected virtual width %d, got %d", expected_width, vs.Width)
	}

	ctrl := vs.GetMachine("controller")
	right := vs.GetMachine("right1")

	if ctrl.OffsetX != 0 {
		t.Errorf("Expected controller OffsetX 0, got %d", ctrl.OffsetX)
	}

	if right.OffsetX != 1920 {
		t.Errorf("Expected right1 OffsetX 1920, got %d", right.OffsetX)
	}

}

func TestComputeVirtualScreen_MultipleLeftOf(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["left1"]      = makeMachine("left1", "left-of", 1280, 720)
	machines["left2"]      = makeMachine("left2", "left-of", 1024, 768)

	vs := computeVirtualScreen("controller", machines)

	expected_width := uint(1024 + 1280 + 1920)

	if vs.Width != expected_width {
		t.Errorf("Expected virtual width %d, got %d", expected_width, vs.Width)
	}

	ctrl := vs.GetMachine("controller")
	left1 := vs.GetMachine("left1")
	left2 := vs.GetMachine("left2")

	if left2.OffsetX != 0 {
		t.Errorf("Expected left2 OffsetX 0, got %d", left2.OffsetX)
	}

	if left1.OffsetX != 1024 {
		t.Errorf("Expected left1 OffsetX 1024, got %d", left1.OffsetX)
	}

	if ctrl.OffsetX != 1024+1280 {
		t.Errorf("Expected controller OffsetX %d, got %d", 1024+1280, ctrl.OffsetX)
	}

}

func TestComputeVirtualScreen_LeftAndRight(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["left1"]      = makeMachine("left1", "left-of", 1280, 720)
	machines["right1"]     = makeMachine("right1", "right-of", 1280, 720)

	vs := computeVirtualScreen("controller", machines)

	expected_width := uint(1280 + 1920 + 1280)

	if vs.Width != expected_width {
		t.Errorf("Expected virtual width %d, got %d", expected_width, vs.Width)
	}

	ctrl  := vs.GetMachine("controller")
	left  := vs.GetMachine("left1")
	right := vs.GetMachine("right1")

	if left.OffsetX != 0 {
		t.Errorf("Expected left1 OffsetX 0, got %d", left.OffsetX)
	}

	if ctrl.OffsetX != 1280 {
		t.Errorf("Expected controller OffsetX 1280, got %d", ctrl.OffsetX)
	}

	if right.OffsetX != 1280+1920 {
		t.Errorf("Expected right1 OffsetX %d, got %d", 1280+1920, right.OffsetX)
	}

}

func TestComputeVirtualScreen_AboveController(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["above1"]     = makeMachine("above1", "above", 1920, 1080)

	vs := computeVirtualScreen("controller", machines)

	if vs.Height != 1080+1080 {
		t.Errorf("Expected virtual height %d, got %d", 1080+1080, vs.Height)
	}

	ctrl  := vs.GetMachine("controller")
	above := vs.GetMachine("above1")

	if above.OffsetY != 0 {
		t.Errorf("Expected above1 OffsetY 0, got %d", above.OffsetY)
	}

	if ctrl.OffsetY != 1080 {
		t.Errorf("Expected controller OffsetY 1080, got %d", ctrl.OffsetY)
	}

}

func TestComputeVirtualScreen_BelowController(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["below1"]     = makeMachine("below1", "below", 1920, 1080)

	vs := computeVirtualScreen("controller", machines)

	if vs.Height != 1080+1080 {
		t.Errorf("Expected virtual height %d, got %d", 1080+1080, vs.Height)
	}

	ctrl  := vs.GetMachine("controller")
	below := vs.GetMachine("below1")

	if ctrl.OffsetY != 0 {
		t.Errorf("Expected controller OffsetY 0, got %d", ctrl.OffsetY)
	}

	if below.OffsetY != 1080 {
		t.Errorf("Expected below1 OffsetY 1080, got %d", below.OffsetY)
	}

}

func TestComputeVirtualScreen_ControllerNotFound(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["left1"] = makeMachine("left1", "left-of", 1280, 720)

	vs := computeVirtualScreen("nonexistent", machines)

	if vs.Width != 1280 {
		t.Errorf("Expected virtual width 1280, got %d", vs.Width)
	}

}

func TestComputeVirtualScreen_MixedPositions(t *testing.T) {

	machines := make(map[string]*Machine)
	machines["controller"] = makeMachine("controller", "center", 1920, 1080)
	machines["left1"]      = makeMachine("left1", "left-of", 1280, 720)
	machines["right1"]     = makeMachine("right1", "right-of", 1280, 720)
	machines["above1"]     = makeMachine("above1", "above", 1920, 1080)
	machines["below1"]     = makeMachine("below1", "below", 1920, 1080)

	vs := computeVirtualScreen("controller", machines)

	expected_width := uint(1280 + 1920 + 1280)

	if vs.Width != expected_width {
		t.Errorf("Expected virtual width %d, got %d", expected_width, vs.Width)
	}

	expected_height := uint(1080 + 1080 + 1080)

	if vs.Height != expected_height {
		t.Errorf("Expected virtual height %d, got %d", expected_height, vs.Height)
	}

	// Verify all machines are in the virtual screen
	for _, hostname := range []string{"controller", "left1", "right1", "above1", "below1"} {
		if vs.GetMachine(hostname) == nil {
			t.Errorf("Expected %s in virtual screen", hostname)
		}
	}

}
