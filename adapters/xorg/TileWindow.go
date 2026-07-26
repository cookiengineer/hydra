package xorg

import "errors"
import "fmt"
import "github.com/cookiengineer/hydra/types"

func TileWindow(bridge *Bridge, window_id uint64, position types.TilePosition, monitor *types.Monitor) error {

	if bridge.display == nil {
		return errors.New("Display is nil")
	}

	if monitor == nil {
		return errors.New("Monitor is nil")
	}

	area_x := monitor.OffsetX
	area_y := monitor.OffsetY
	area_w := monitor.Width
	area_h := monitor.Height

	var target_x, target_y, target_w, target_h int

	switch position {
	case types.TileLeft:
		target_x = area_x
		target_y = area_y
		target_w = area_w / 2
		target_h = area_h

	case types.TileRight:
		target_x = area_x + area_w/2
		target_y = area_y
		target_w = area_w / 2
		target_h = area_h

	case types.TileTop:
		target_x = area_x
		target_y = area_y
		target_w = area_w
		target_h = area_h / 2

	case types.TileBottom:
		target_x = area_x
		target_y = area_y + area_h/2
		target_w = area_w
		target_h = area_h / 2

	case types.TileTopLeft:
		target_x = area_x
		target_y = area_y
		target_w = area_w / 2
		target_h = area_h / 2

	case types.TileTopRight:
		target_x = area_x + area_w/2
		target_y = area_y
		target_w = area_w / 2
		target_h = area_h / 2

	case types.TileBottomLeft:
		target_x = area_x
		target_y = area_y + area_h/2
		target_w = area_w / 2
		target_h = area_h / 2

	case types.TileBottomRight:
		target_x = area_x + area_w/2
		target_y = area_y + area_h/2
		target_w = area_w / 2
		target_h = area_h / 2

	default:
		return fmt.Errorf("Unknown tile position: %s", position.String())
	}

	err0 := MoveResizeWindow(bridge, window_id, target_x, target_y, target_w, target_h)

	if err0 != nil {
		return err0
	}

	err1 := RaiseWindow(bridge, window_id)

	if err1 != nil {
		return err1
	}

	return nil

}
