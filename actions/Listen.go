package actions

import "context"
import "encoding/json"
import "fmt"
import "os"
import "os/signal"
import "syscall"
import "github.com/cookiengineer/hydra/listeners"

func Listen(host string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel,
		os.Interrupt,    // Ctrl+C (SIGINT)
		syscall.SIGTERM, // kill <pid>
		syscall.SIGHUP,
	)

	go func() {

		<-signal_channel
		cancel()

	}()

	state, err0 := listeners.Init(":0")

	if err0 == nil {

		go listeners.StartLoop(state)

		go func() {

			for {
				select {
				case event := <-state.MouseEvents:

					data, _ := json.Marshal(event)
					fmt.Printf("Mouse: %+v\n", string(data))

					// TODO: send to correct client

				case event  := <-state.KeyboardEvents:

					data, _ := json.Marshal(event)
					fmt.Printf("Key: %+v\n", string(data))

					// TODO: send to correct client

				case <-ctx.Done():

					return

				}
			}

		}()

		<-ctx.Done()

		fmt.Println("Shutting down...")

		state.Destroy()

		return nil

	} else {
		return err0
	}

}

