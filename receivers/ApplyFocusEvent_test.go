package receivers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestApplyFocusEvent_NoopWithNilBridge(t *testing.T) {

	event := &types.FocusEvent{
		Type:      "focus",
		Direction: "left",
	}

	ApplyFocusEvent(nil, event)

}

func TestApplyFocusEvent_AllDirections(t *testing.T) {

	directions := []string{"left", "right", "up", "down"}

	for _, dir := range directions {

		event := &types.FocusEvent{
			Type:      "focus",
			Direction: dir,
		}

		ApplyFocusEvent(nil, event)

	}

}
