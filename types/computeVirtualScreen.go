package types

func computeVirtualScreen(controller string, machines map[string]*Machine) *Screen {

	virtual_min_x := uint(0)
	virtual_min_y := uint(0)
	virtual_max_x := uint(0)
	virtual_max_y := uint(0)

	controller_machine, ok := machines[controller]

	if ok == true {

		virtual_min_x = controller_machine.Screen.OffsetX
		virtual_min_y = controller_machine.Screen.OffsetY
		virtual_max_x = controller_machine.Screen.OffsetX + controller_machine.Screen.Width
		virtual_max_y = controller_machine.Screen.OffsetY + controller_machine.Screen.Height

	}


	left_to_right := make([]*Machine, 0)
	top_to_bottom := make([]*Machine, 0)

	for _, machine := range machines {

		switch machine.Position {
		case "left-of":

			left_to_right = make([]*Machine, 0)
			left_to_right = append(left_to_right, machine)
			left_to_right = append(left_to_right, left_to_right...)

		case "right-of":

			left_to_right = append(left_to_right, machine)

		case "center":

			left  := make([]*Machine, 0)
			right := make([]*Machine, 0)

			for _, other := range left_to_right {

				if other.Position == "left-of" {
					left = append(left, other)
				} else if other.Position == "right-of" {
					right = append(right, other)
				}

			}

			left_to_right = make([]*Machine, 0)
			left_to_right = append(left_to_right, left...)
			left_to_right = append(left_to_right, machine)
			left_to_right = append(left_to_right, right...)

		}

	}

	for _, machine := range machines {

		switch machine.Position {
		case "above":

			top_to_bottom = make([]*Machine, 0)
			top_to_bottom = append(top_to_bottom, machine)
			top_to_bottom = append(top_to_bottom, top_to_bottom...)

		case "below":

			top_to_bottom = append(top_to_bottom, machine)

		case "center":

			above := make([]*Machine, 0)
			below := make([]*Machine, 0)

			for _, other := range top_to_bottom {

				if other.Position == "above" {
					above = append(above, other)
				} else if other.Position == "below" {
					below = append(below, other)
				}

			}

			top_to_bottom = make([]*Machine, 0)
			top_to_bottom = append(top_to_bottom, above...)
			top_to_bottom = append(top_to_bottom, machine)
			top_to_bottom = append(top_to_bottom, below...)

		}

	}

	relative_offset_x := uint(0)
	relative_offset_y := uint(0)

	for _, machine := range left_to_right {

		machine.Screen.OffsetX = relative_offset_x

		if machine.Screen.OffsetX + machine.Screen.Width > virtual_max_x {
			virtual_max_x = machine.Screen.OffsetX + machine.Screen.Width
		}

		if machine.Screen.Height > virtual_max_y {
			virtual_max_y = machine.Screen.Height
		}

		relative_offset_x += machine.Screen.Width

	}

	for _, machine := range top_to_bottom {

		machine.Screen.OffsetY = relative_offset_y

		if machine.Screen.OffsetY + machine.Screen.Height > virtual_max_y {
			virtual_max_y = machine.Screen.OffsetY + machine.Screen.Height
		}

		if machine.Screen.Width > virtual_max_x {
			virtual_max_x = machine.Screen.Width
		}

		relative_offset_y += machine.Screen.Height

	}


	return &Screen{
		Width:    virtual_max_x,
		Height:   virtual_max_y,
		Monitors: []Monitor{},
		OffsetX:  virtual_min_x,
		OffsetY:  virtual_min_y,
	}

}
