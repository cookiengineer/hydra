package types

type TilePosition int

const (
	TileLeft        TilePosition = 0
	TileRight       TilePosition = 1
	TileTopLeft     TilePosition = 2
	TileTopRight    TilePosition = 3
	TileBottomLeft  TilePosition = 4
	TileBottomRight TilePosition = 5
)

func (position TilePosition) String() string {

	switch position {
	case TileLeft:
		return "left"
	case TileRight:
		return "right"
	case TileTopLeft:
		return "top-left"
	case TileTopRight:
		return "top-right"
	case TileBottomLeft:
		return "bottom-left"
	case TileBottomRight:
		return "bottom-right"
	default:
		return "unknown"
	}

}

func TilePositionFromString(tile  string) TilePosition {

	switch tile {
	case "left":
		return TileLeft
	case "right":
		return TileRight
	case "top-left":
		return TileTopLeft
	case "top-right":
		return TileTopRight
	case "bottom-left":
		return TileBottomLeft
	case "bottom-right":
		return TileBottomRight
	default:
		return -1
	}

}
