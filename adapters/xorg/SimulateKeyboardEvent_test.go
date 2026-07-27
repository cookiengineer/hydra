//go:build integration

package xorg

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestSimulateKeyboardEvent_Press(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateKeyPress(bridge, 10)

	if err != nil {
		t.Errorf("SimulateKeyPress failed: %v", err)
	}

}

func TestSimulateKeyboardEvent_Release(t *testing.T) {

	bridge := createTestBridge(t)

	err := SimulateKeyRelease(bridge, 20)

	if err != nil {
		t.Errorf("SimulateKeyRelease failed: %v", err)
	}

}

func TestSimulateKeyboardEvent_DispatchPress(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.KeyboardEvent{
		Type:    types.KeyPress,
		Keycode: 33,
	}

	SimulateKeyboardEvent(bridge, event)

}

func TestSimulateKeyboardEvent_DispatchRelease(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.KeyboardEvent{
		Type:    types.KeyRelease,
		Keycode: 44,
	}

	SimulateKeyboardEvent(bridge, event)

}

func TestSimulateKeyboardEvent_NilBridge(t *testing.T) {

	bridge := createTestBridge(t)

	event := &types.KeyboardEvent{
		Type:    types.KeyPress,
		Keycode: 42,
	}

	SimulateKeyboardEvent(bridge, event)

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SimulateKeyboardEvent panicked with nil bridge: %v", r)
			}
		}()

		var nil_bridge *Bridge
		SimulateKeyboardEvent(nil_bridge, event)
	}()

}
