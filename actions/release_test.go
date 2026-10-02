package actions

import "encoding/json"
import "testing"
import "github.com/cookiengineer/hydra/types"

func TestReleaseTrackedInput(t *testing.T) {

	state := types.NewGlobalState()
	state.TrackKey(0xFFEB)
	state.TrackKey(0x61)
	state.TrackButton(1)

	machine := &types.Machine{
		Hostname: "hydratwo",
		Socket:   make(chan []byte, 8),
	}

	releaseTrackedInput(state, machine, 100, 200)

	if len(state.TrackedKeys()) != 0 {
		t.Errorf("Expected tracked keys to be cleared, got %d", len(state.TrackedKeys()))
	}

	if len(state.TrackedButtons()) != 0 {
		t.Errorf("Expected tracked buttons to be cleared, got %d", len(state.TrackedButtons()))
	}

	types_seen := make([]string, 0)

drain:
	for {
		select {
		case data := <-machine.Socket:
			var envelope types.Event
			if err := json.Unmarshal(data, &envelope); err != nil {
				t.Fatalf("Failed to unmarshal envelope: %v", err)
			}
			types_seen = append(types_seen, envelope.Type)
		default:
			break drain
		}
	}

	if len(types_seen) != 3 {
		t.Fatalf("Expected 3 release events, got %d (%v)", len(types_seen), types_seen)
	}

	keyboards := 0
	mice := 0

	for _, event_type := range types_seen {
		if event_type == "keyboard" {
			keyboards++
		}
		if event_type == "mouse" {
			mice++
		}
	}

	if keyboards != 2 {
		t.Errorf("Expected 2 keyboard releases, got %d", keyboards)
	}

	if mice != 1 {
		t.Errorf("Expected 1 mouse release, got %d", mice)
	}

}

func TestReleaseTrackedInputNilSocket(t *testing.T) {

	state := types.NewGlobalState()
	state.TrackKey(0xFFE1)
	state.TrackButton(3)

	releaseTrackedInput(state, &types.Machine{Hostname: "gone"}, 0, 0)

	if len(state.TrackedKeys()) != 0 {
		t.Errorf("Expected tracked keys to be cleared")
	}

	if len(state.TrackedButtons()) != 0 {
		t.Errorf("Expected tracked buttons to be cleared")
	}

}
