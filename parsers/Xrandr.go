package parsers

import (
	"bufio"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cookiengineer/hydra/types"
)

func Xrandr() (*types.Screen, error) {

	cmd := exec.Command("xrandr", "--query")

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	screen := &types.Screen{}

	var currentMonitor *types.Monitor

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {

		line := scanner.Text()

		// Parse global virtual screen size
		// Example:
		// Screen 0: minimum 8 x 8, current 3840 x 1080, maximum 32767 x 32767
		if strings.HasPrefix(line, "Screen ") {

			parts := strings.Split(line, "current")

			if len(parts) > 1 {

				right := strings.TrimSpace(parts[1])
				right = strings.Split(right, ",")[0]
				dims := strings.Split(strings.TrimSpace(right), " x ")

				if len(dims) == 2 {
					w, _ := strconv.Atoi(strings.TrimSpace(dims[0]))
					h, _ := strconv.Atoi(strings.TrimSpace(dims[1]))
					screen.Width = w
					screen.Height = h
				}
			}
			continue
		}

		// Parse connected monitors
		// Example:
		// HDMI-1 connected primary 1920x1080+0+0 ...
		if strings.Contains(line, " connected") {

			fields := strings.Fields(line)

			if len(fields) < 3 {
				continue
			}

			monitor := types.Monitor{
				Output:    fields[0],
				Connected: true,
			}

			for _, f := range fields {

				// resolution+offset pattern 1920x1080+0+0
				if strings.Contains(f, "x") && strings.Contains(f, "+") {

					resOffset := strings.Split(f, "+")
					if len(resOffset) >= 3 {

						monitor.Resolution = resOffset[0]

						dims := strings.Split(resOffset[0], "x")
						if len(dims) == 2 {
							monitor.Width, _ = strconv.Atoi(dims[0])
							monitor.Height, _ = strconv.Atoi(dims[1])
						}

						monitor.OffsetX, _ = strconv.Atoi(resOffset[1])
						monitor.OffsetY, _ = strconv.Atoi(resOffset[2])
					}
				}
			}

			screen.Monitors = append(screen.Monitors, monitor)
			currentMonitor = &screen.Monitors[len(screen.Monitors)-1]

			continue
		}

		// Parse modes (must belong to last connected monitor)
		// Example:
		//   1920x1080     60.00*+  59.94
		if currentMonitor != nil {

			trimmed := strings.TrimSpace(line)

			if strings.Contains(trimmed, "x") {

				fields := strings.Fields(trimmed)
				if len(fields) >= 2 {

					mode := types.MonitorMode{
						Resolution: fields[0],
					}

					refreshStr := fields[1]
					refreshStr = strings.TrimRight(refreshStr, "*+")
					refresh, err := strconv.ParseFloat(refreshStr, 32)
					if err == nil {
						mode.RefreshRate = float32(refresh)
					}

					currentMonitor.Modes = append(currentMonitor.Modes, mode)
				}
			}
		}
	}

	if screen.Width == 0 || screen.Height == 0 {
		return nil, errors.New("could not parse xrandr screen size")
	}

	return screen, nil
}
