package main

import "github.com/cookiengineer/hydra/actions"
import "fmt"
import "os"
import "os/user"
import "strings"

func showUsage() {

	fmt.Println("Usage: ")
	fmt.Println("  hydra listen <host>")
	fmt.Println("  hydra connect left-of <host>")
	fmt.Println("  hydra connect right-of <host>")
	fmt.Println("  hydra connect above <host>")
	fmt.Println("  hydra connect below <host>")

}

func main() {

	display := os.Getenv("DISPLAY")
	action  := ""
	host    := ""
	position := ""

	if display == "" {
		display = ":0"
	}

	current_user, err0 := user.Current()

	if len(os.Args) == 3 {

		if os.Args[1] == "listen" {
			action = "listen"
			host   = strings.TrimSpace(strings.ToLower(os.Args[2]))
		}

	} else if len(os.Args) == 4 {

		if os.Args[1] == "connect" {
			action   = "connect"
			position = strings.TrimSpace(strings.ToLower(os.Args[2]))
			host     = strings.TrimSpace(strings.ToLower(os.Args[3]))
		}

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

		} else if action == "connect" {

			err1 := actions.Connect(host, position)

			if err1 != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err1.Error())
				os.Exit(1)
			}

			err2 := actions.ReceiveEvents(host)

			if err2 != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err2.Error())
				os.Exit(1)
			}

			os.Exit(0)

		} else {
			showUsage()
			os.Exit(1)
		}

	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err0.Error())
		os.Exit(1)
	}

}
