//go:build integration

package xorg

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestSimulateMouseEvent_Move(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.MouseEvent{
		Type: types.MouseMove,
		X:    100,
		Y:    100,
		DX:   10,
		DY:   20,
	}

	SimulateMouseEvent(bridge, event)

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 110 || y != 120 {
		t.Errorf("Expected pointer at (110, 120), got (%d, %d)", x, y)
	}

}

func TestSimulateMouseEvent_Press(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.MouseEvent{
		Type:   types.MouseButtonPress,
		X:      200,
		Y:      200,
		Button: types.MouseButtonRight,
	}

	SimulateMouseEvent(bridge, event)

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 200 || y != 200 {
		t.Errorf("Expected pointer at (200, 200), got (%d, %d)", x, y)
	}

}

func TestSimulateMouseEvent_Release(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.MouseEvent{
		Type:   types.MouseButtonRelease,
		X:      300,
		Y:      300,
		Button: types.MouseButtonMiddle,
	}

	SimulateMouseEvent(bridge, event)

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 300 || y != 300 {
		t.Errorf("Expected pointer at (300, 300), got (%d, %d)", x, y)
	}

}

func TestSimulateMouseEvent_Scroll(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.MouseEvent{
		Type: types.MouseScroll,
		X:    400,
		Y:    400,
		DX:   0,
		DY:   2,
	}

	SimulateMouseEvent(bridge, event)

	x, y, err := bridge.QueryPointer()

	if err != nil {
		t.Errorf("QueryPointer failed: %v", err)
		return
	}

	if x != 400 || y != 400 {
		t.Errorf("Expected pointer at (400, 400), got (%d, %d)", x, y)
	}

}

func TestSimulateMouseEvent_NilBridge(t *testing.T) {

	event := &types.MouseEvent{
		Type: types.MouseMove,
		X:    100,
		Y:    100,
		DX:   0,
		DY:   0,
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SimulateMouseEvent panicked with nil bridge: %v", r)
			}
		}()

		var nil_bridge *Bridge
		SimulateMouseEvent(nil_bridge, event)
	}()

}
