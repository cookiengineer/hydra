package xorg

import "fmt"
import "testing"

func TestBridge(t *testing.T) {

	t.Run("SimulateMouseMove()", func(t *testing.T) {

		bridge, err0 := NewBridge(":0")

		if err0 != nil {
			t.Errorf("Expected %v to be nil", err0)
		}

		err1 := SimulateMouseMove(bridge, 1300, 100, 37, 37)

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		x, y, err2 := bridge.QueryPointer()

		if err2 == nil {

			if x != 1337 || y != 137 {
				t.Errorf("Expected %d x %d to be %d x %d", x, y, 1337, 137)
			}

		} else {
			t.Errorf("Expected %v to be nil", err2)
		}

		fmt.Printf("Mouse Events length: %d\n", len(bridge.MouseEvents))

	})

}
