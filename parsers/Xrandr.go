package parsers

import "fmt"
import "os/exec"
import "regexp"
import "strconv"
import "strings"
import "github.com/cookiengineer/hydra/types"

func Xrandr() (*types.Screen, error) {

	cmd := exec.Command("xrandr", "--query")
	buffer, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return parseXrandrOutput(string(buffer))

}

func parseXrandrOutput(output string) (*types.Screen, error) {

	lines := strings.Split(strings.TrimSpace(output), "\n")
	screen := &types.Screen{}

	for _, line := range lines {

		if strings.HasPrefix(line, "Screen ") && strings.Contains(line, ": ") && strings.Contains(line, " current ") {

			pattern := regexp.MustCompile(`current\s+([0-9]+)\s+x\s+([0-9]+)`)
			matches := pattern.FindStringSubmatch(line)

			if len(matches) == 3 {

				width, err1 := strconv.Atoi(matches[1])
				height, err2 := strconv.Atoi(matches[2])

				if err1 == nil && err2 == nil {
					screen.Width = uint(width)
					screen.Height = uint(height)
				}

			}

		} else if strings.Contains(line, " connected ") {

			pattern := regexp.MustCompile(`^(\S+)\s+connected(?:\s+primary)?\s+(\d+)x(\d+)\+(\d+)\+(\d+)`)
			matches := pattern.FindStringSubmatch(line)

			if len(matches) == 6 {

				width, err1 := strconv.Atoi(matches[2])
				height, err2 := strconv.Atoi(matches[3])
				offset_x, err3 := strconv.Atoi(matches[4])
				offset_y, err4 := strconv.Atoi(matches[5])

				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {

					screen.Monitors = append(screen.Monitors, types.Monitor{
						Output:     matches[1],
						Connected:  true,
						Resolution: fmt.Sprintf("%dx%d", width, height),
						Width:      width,
						Height:     height,
						OffsetX:    offset_x,
						OffsetY:    offset_y,
					})

				}

			}

		}

	}

	return screen, nil

}
