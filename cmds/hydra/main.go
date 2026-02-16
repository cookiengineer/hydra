package main

import "github.com/cookiengineer/hydra/listeners"
import "github.com/cookiengineer/hydra/types"
import "encoding/json"
import "fmt"
import "os"
import "os/user"

func main() {

	events   := make(chan types.MouseEvent, 128)
	xdisplay := os.Getenv("DISPLAY")

	if xdisplay == "" {
		xdisplay = ":0"
	}

	current_user, err0 := user.Current()

	if err0 == nil {

		fmt.Println("USER=" + current_user.Username)
		fmt.Println("DISPLAY=" + xdisplay)

		err1 := listeners.CaptureMouse(events, xdisplay)

		if err1 == nil {
			fmt.Println("Listening ...")
		} else {
			panic(err1)
		}
	
	} else {
		panic(err0)
	}

	for event := range events {
		data, _ := json.Marshal(event)
		fmt.Printf("%+v\n", string(data))
	}

}
