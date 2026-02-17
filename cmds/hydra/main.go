package main

import "github.com/cookiengineer/hydra/actions"
import "fmt"
import "os"
import "os/user"
import "strings"

func showUsage() {

	fmt.Println("Usage: ")
	fmt.Println("  hydra listen <host>")

}

func main() {

	display := os.Getenv("DISPLAY")
	action  := ""
	host    := ""

	if display == "" {
		display = ":0"
	}

	current_user, err0 := user.Current()

	if len(os.Args) == 3 {

		if os.Args[1] == "listen" {
			action = "listen"
			host   = strings.TrimSpace(strings.ToLower(os.Args[2]))
		}

	} else {
		action = ""
	}

	if err0 == nil {

		fmt.Println("USER=" + current_user.Username)
		fmt.Println("DISPLAY=" + display)

		if action == "listen" {

			err1 := actions.Listen(host)

			if err1 != nil {

				fmt.Fprintf(os.Stderr, "Error: %s\n", err1.Error())
				os.Exit(1)

			} else {
				os.Exit(0)
			}

		} else {
			showUsage()
			os.Exit(1)
		}
	
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err0.Error())
		os.Exit(1)
	}

}
