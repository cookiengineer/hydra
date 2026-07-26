package types

type TilePosition int

const (
	TileLeft        TilePosition = 0
	TileRight       TilePosition = 1
	TileTop         TilePosition = 2
	TileBottom      TilePosition = 3
	TileTopLeft     TilePosition = 4
	TileTopRight    TilePosition = 5
	TileBottomLeft  TilePosition = 6
	TileBottomRight TilePosition = 7
)

func (position TilePosition) String() string {

	switch position {
	case TileLeft:
		return "left"
	case TileRight:
		return "right"
	case TileTop:
		return "top"
	case TileBottom:
		return "bottom"
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

func TilePositionFromString(tile string) TilePosition {

	switch tile {
	case "left":
		return TileLeft
	case "right":
		return TileRight
	case "top":
		return TileTop
	case "bottom":
		return TileBottom
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
