package receivers

import "testing"
import "github.com/cookiengineer/hydra/types"

func resetTracked() {

	pressedMutex.Lock()
	defer pressedMutex.Unlock()

	pressedKeys = make(map[uint32]bool)
	pressedButtons = make(map[types.MouseEventButton]bool)

}

func TestTrackKeyAndButton(t *testing.T) {

	resetTracked()

	trackKey(114, true)
	trackButton(types.MouseButtonLeft, true)

	pressedMutex.Lock()

	if !pressedKeys[114] {
		t.Error("Expected key 114 to be tracked")
	}

	if !pressedButtons[types.MouseButtonLeft] {
		t.Error("Expected left button to be tracked")
	}

	pressedMutex.Unlock()

	trackKey(114, false)
	trackButton(types.MouseButtonLeft, false)

	pressedMutex.Lock()

	if pressedKeys[114] {
		t.Error("Expected key 114 to be untracked")
	}

	if pressedButtons[types.MouseButtonLeft] {
		t.Error("Expected left button to be untracked")
	}

	pressedMutex.Unlock()

	resetTracked()

}

func TestTrackKeyZeroIgnored(t *testing.T) {

	resetTracked()

	trackKey(0, true)

	pressedMutex.Lock()

	if len(pressedKeys) != 0 {
		t.Errorf("Expected keycode 0 to be ignored, got %d tracked", len(pressedKeys))
	}

	pressedMutex.Unlock()

	resetTracked()

}
