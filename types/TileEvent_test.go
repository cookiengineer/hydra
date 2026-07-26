package types

import "encoding/json"
import "testing"

func TestTileEventJSON(t *testing.T) {

	tests := []string{"left", "right", "top", "bottom", "top-left", "top-right", "bottom-left", "bottom-right"}

	for _, position := range tests {

		event := TileEvent{
			Type:     "tile",
			Position: position,
		}

		data, err := json.Marshal(event)

		if err != nil {
			t.Errorf("Failed to marshal TileEvent: %s", err.Error())
		}

		var parsed TileEvent

		err = json.Unmarshal(data, &parsed)

		if err != nil {
			t.Errorf("Failed to unmarshal TileEvent: %s", err.Error())
		}

		if parsed.Type != "tile" {
			t.Errorf("Expected type tile, got %s", parsed.Type)
		}

		if parsed.Position != position {
			t.Errorf("Expected position %s, got %s", position, parsed.Position)
		}

	}

}
