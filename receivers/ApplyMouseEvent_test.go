package receivers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestCoordinateTranslation(t *testing.T) {

	vs := types.NewVirtualScreen()
	vs.AddMachine("client1", &types.Screen{
		Width:   1280,
		Height:  720,
		OffsetX: 1920,
		OffsetY: 0,
	})

	event := &types.MouseEvent{
		Type: types.MouseMove,
		X:    2500,
		Y:    100,
		DX:   5,
		DY:   0,
	}

	ApplyMouseEvent(nil, event, vs, "client1")

	expected_x := uint(580)
	expected_y := uint(100)

	if event.X != expected_x {
		t.Errorf("Expected translated X %d, got %d", expected_x, event.X)
	}

	if event.Y != expected_y {
		t.Errorf("Expected translated Y %d, got %d", expected_y, event.Y)
	}

}

func TestCoordinateTranslation_NoOffset(t *testing.T) {

	vs := types.NewVirtualScreen()
	vs.AddMachine("client1", &types.Screen{
		Width:   1920,
		Height:  1080,
		OffsetX: 0,
		OffsetY: 0,
	})

	event := &types.MouseEvent{
		Type: types.MouseMove,
		X:    500,
		Y:    300,
		DX:   10,
		DY:   5,
	}

	ApplyMouseEvent(nil, event, vs, "client1")

	if event.X != 500 {
		t.Errorf("Expected X 500, got %d", event.X)
	}

	if event.Y != 300 {
		t.Errorf("Expected Y 300, got %d", event.Y)
	}

}

func TestCoordinateTranslation_NegativeOffset(t *testing.T) {

	vs := types.NewVirtualScreen()
	vs.AddMachine("client1", &types.Screen{
		Width:   1920,
		Height:  1080,
		OffsetX: 500,
		OffsetY: 200,
	})

	event := &types.MouseEvent{
		Type: types.MouseMove,
		X:    100,
		Y:    150,
		DX:   0,
		DY:   0,
	}

	ApplyMouseEvent(nil, event, vs, "client1")

	// Since X < OffsetX, no subtraction happens
	if event.X != 100 {
		t.Errorf("Expected X 100, got %d", event.X)
	}

	if event.Y != 150 {
		t.Errorf("Expected Y 150, got %d", event.Y)
	}

}

