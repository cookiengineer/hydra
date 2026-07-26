package receivers

import "testing"
import "github.com/cookiengineer/hydra/types"

func TestApplyTileEvent_NoopWithNilBridge(t *testing.T) {

	vs := types.NewVirtualScreen()
	event := &types.TileEvent{
		Type:     "tile",
		Position: "left",
	}

	ApplyTileEvent(nil, event, vs, "client1")

}

func TestApplyTileEvent_NoopWithNilVirtualScreen(t *testing.T) {

	event := &types.TileEvent{
		Type:     "tile",
		Position: "right",
	}

	ApplyTileEvent(nil, event, nil, "client1")

}

func TestApplyTileEvent_AllPositions(t *testing.T) {

	vs := types.NewVirtualScreen()
	positions := []string{"left", "right", "top", "bottom", "top-left", "top-right", "bottom-left", "bottom-right"}

	for _, pos := range positions {

		event := &types.TileEvent{
			Type:     "tile",
			Position: pos,
		}

		ApplyTileEvent(nil, event, vs, "client1")

	}

}
