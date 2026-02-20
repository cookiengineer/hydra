package parsers

import "fmt"
import "testing"

func TestXrandr(t *testing.T) {

	t.Run("Xrandr()", func(t *testing.T) {

		screen, err := Xrandr()

		if err == nil {

			if screen.Width != 3840 {
				t.Errorf("Expected screen.Width %d to be %d", screen.Width, 3840)
			}

			if screen.Height != 1080 {
				t.Errorf("Expected screen.Height %d to be %d", screen.Height, 1080)
			}

			if len(screen.Monitors) == 2 {

				if screen.Monitors[0].Output != "HDMI-A-0" {
					t.Errorf("Expected %s to be %s", screen.Monitors[0].Output, "HDMI-A-0")
				}

				if screen.Monitors[0].Connected != true {
					t.Errorf("Expected %v to be %v", screen.Monitors[0].Connected, true)
				}

				if screen.Monitors[0].Resolution != "1920x1080" {
					t.Errorf("Expected %s to be %s", screen.Monitors[0].Resolution, "1920x1080")
				}

				if screen.Monitors[0].Width != 1920 {
					t.Errorf("Expected %d to be %d", screen.Monitors[0].Width, 1920)
				}

				if screen.Monitors[0].Height != 1080 {
					t.Errorf("Expected %d to be %d", screen.Monitors[0].Height, 1080)
				}

			} else {
				t.Errorf("Expected %d screen.Monitors to be %d", len(screen.Monitors), 2)
			}

			fmt.Println(screen)

		} else {
			t.Errorf("Expected %v to be nil", err)
		}

	})

}
