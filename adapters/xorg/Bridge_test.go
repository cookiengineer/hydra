//go:build integration

package xorg

import "fmt"
import "os"
import "testing"

func TestBridge(t *testing.T) {

	display := os.Getenv("DISPLAY")

	if display == "" {
		t.Skip("DISPLAY not set; start Xvfb: Xvfb :99 -screen 0 1920x1080x24 &")
	}

	t.Run("SimulateMouseMove()", func(t *testing.T) {

		bridge, err0 := NewBridge(display)

		if err0 != nil {
			t.Fatalf("Expected %v to be nil", err0)
		}

		t.Cleanup(func() {
			bridge.Destroy()
		})

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
