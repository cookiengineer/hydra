package actions

import "os"

func defaultDisplay() string {

	display := os.Getenv("DISPLAY")

	if display == "" {
		display = ":0"
	}

	return display

}
