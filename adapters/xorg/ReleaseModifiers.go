package xorg

func ReleaseModifiers(bridge *Bridge) {

	if bridge == nil || bridge.display == nil {
		return
	}

	keysyms := []uint32{
		0xFFE1, // Shift_L
		0xFFE2, // Shift_R
		0xFFE3, // Control_L
		0xFFE4, // Control_R
		0xFFE7, // Meta_L
		0xFFE8, // Meta_R
		0xFFE9, // Alt_L
		0xFFEA, // Alt_R
		0xFFEB, // Super_L
		0xFFEC, // Super_R
		0xFE03, // ISO_Level3_Shift
		0xFF7E, // Mode_switch
	}

	for _, keysym := range keysyms {

		keycode, err := KeysymToKeycode(bridge, keysym)

		if err == nil && keycode != 0 {
			SimulateKeyRelease(bridge, keycode)
		}

	}

}
