package types

func computeVirtualScreen(host Machine, machines []Machine) *Screen {

	screen := Screen{}

	// Start with host dimensions
	hostWidth := host.Screen.Width
	hostHeight := host.Screen.Height

	minX := 0
	minY := 0
	maxX := hostWidth
	maxY := hostHeight

	for _, machine := range machines {

		w := machine.Screen.Width
		h := machine.Screen.Height

		switch machine.Position {

		case "right-of":
			maxX += w
			if h > maxY {
				maxY = h
			}

		case "left-of":
			minX -= w
			if h > maxY {
				maxY = h
			}

		case "below":
			maxY += h
			if w > maxX {
				maxX = w
			}

		case "above":
			minY -= h
			if w > maxX {
				maxX = w
			}
		}
	}

	screen.Width = maxX - minX
	screen.Height = maxY - minY

	// Recompute monitor offsets into unified coordinate space

	offsetXShift := -minX
	offsetYShift := -minY

	// Host monitors
	for _, m := range host.Screen.Monitors {

		monitor := m
		monitor.OffsetX += offsetXShift
		monitor.OffsetY += offsetYShift

		screen.Monitors = append(screen.Monitors, monitor)

	}

	// Remote monitors
	currentRightX := hostWidth
	currentLeftX := 0
	currentTopY := 0
	currentBottomY := hostHeight

	for _, machine := range machines {

		for _, m := range machine.Screen.Monitors {

			monitor := m

			switch machine.Position {

			case "right-of":
				monitor.OffsetX += currentRightX + offsetXShift
				monitor.OffsetY += offsetYShift
			case "left-of":
				monitor.OffsetX += currentLeftX - machine.Screen.Width + offsetXShift
				monitor.OffsetY += offsetYShift
			case "below":
				monitor.OffsetX += offsetXShift
				monitor.OffsetY += currentBottomY + offsetYShift
			case "above":
				monitor.OffsetX += offsetXShift
				monitor.OffsetY += currentTopY - machine.Screen.Height + offsetYShift
			}

			screen.Monitors = append(screen.Monitors, monitor)

		}

	}

	return &screen

}
