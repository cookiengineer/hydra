package parsers

import "testing"

func TestParseXrandrOutput(t *testing.T) {

	t.Run("SingleMonitor", func(t *testing.T) {

		input := `Screen 0: minimum 320 x 200, current 1920 x 1080, maximum 16384 x 16384
HDMI-A-0 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 531mm x 299mm
   1920x1080     60.00*+  50.00    59.94
   1680x1050     59.95
   1280x1024     75.02    60.02
   1440x900      74.98    59.89
   1280x960      60.00
   1152x864      75.00
   1024x768      75.03    70.07    60.00
   800x600       75.00    72.19    60.32    56.25
   640x480       75.00    72.81    59.94
DP-0 disconnected (normal left inverted right x axis y axis)
DP-1 disconnected (normal left inverted right x axis y axis)
`

		screen, err := parseXrandrOutput(input)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if screen.Width != 1920 {
			t.Errorf("Expected width 1920, got %d", screen.Width)
		}

		if screen.Height != 1080 {
			t.Errorf("Expected height 1080, got %d", screen.Height)
		}

		if len(screen.Monitors) != 1 {
			t.Errorf("Expected 1 monitor, got %d", len(screen.Monitors))
		}

		if screen.Monitors[0].Output != "HDMI-A-0" {
			t.Errorf("Expected HDMI-A-0, got %s", screen.Monitors[0].Output)
		}

		if !screen.Monitors[0].Connected {
			t.Error("Expected connected")
		}

		if screen.Monitors[0].Width != 1920 {
			t.Errorf("Expected monitor width 1920, got %d", screen.Monitors[0].Width)
		}

	})

	t.Run("DualMonitor", func(t *testing.T) {

		input := `Screen 0: minimum 320 x 200, current 3840 x 1080, maximum 16384 x 16384
HDMI-A-0 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 531mm x 299mm
HDMI-A-1 connected 1920x1080+1920+0 (normal left inverted right x axis y axis) 531mm x 299mm
DP-0 disconnected (normal left inverted right x axis y axis)
`

		screen, err := parseXrandrOutput(input)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if screen.Width != 3840 {
			t.Errorf("Expected width 3840, got %d", screen.Width)
		}

		if screen.Height != 1080 {
			t.Errorf("Expected height 1080, got %d", screen.Height)
		}

		if len(screen.Monitors) != 2 {
			t.Errorf("Expected 2 monitors, got %d", len(screen.Monitors))
		}

		if screen.Monitors[1].OffsetX != 1920 {
			t.Errorf("Expected second monitor OffsetX 1920, got %d", screen.Monitors[1].OffsetX)
		}

	})

	t.Run("NoMonitors", func(t *testing.T) {

		input := `Screen 0: minimum 320 x 200, current 1024 x 768, maximum 16384 x 16384
DP-0 disconnected (normal left inverted right x axis y axis)
DP-1 disconnected (normal left inverted right x axis y axis)
`

		screen, err := parseXrandrOutput(input)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(screen.Monitors) != 0 {
			t.Errorf("Expected 0 monitors, got %d", len(screen.Monitors))
		}

	})

}
