package actions

import "context"
import "encoding/json"
import "fmt"
import "net/http"
import "os"
import "os/signal"
import "sync"
import "syscall"
import "github.com/cookiengineer/hydra/adapters/xorg"
import "github.com/cookiengineer/hydra/handlers"
import "github.com/cookiengineer/hydra/parsers"
import "github.com/cookiengineer/hydra/types"

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

	bridge, err0 := xorg.NewBridge(":0")
	screen, err1 := parsers.Xrandr()

	if err0 == nil && err1 == nil {

		config := types.NewConfig(types.Machine{
			Hostname: host,
			IP:       "", // populated later
			Position: "center",
			Screen:   *screen,
		})

		// Controller is the current machine
		config.SetThis(host)

		http.HandleFunc("/config",    handlers.OnConfig(config))
		// func(response http.ResponseWriter, request *http.Request) {

		// 	config.Lock()
		// 	json.NewEncoder(response).Encode(config)
		// 	config.Unlock()

		// })

		http.HandleFunc("/connect", func(response http.ResponseWriter, request *http.Request) {
			handlers.OnConnect(config, response, request)
		})

		http.HandleFunc("/disconnect", func(response http.ResponseWriter, request *http.Request) {
			handlers.OnDisconnect(config, response, request)
		})

		go http.ListenAndServe(":3000", nil)

		go bridge.Init()

		go func() {

			for {
				select {
				case <-bridge.MouseEvents:

					mouse_x, mouse_y, err0 := bridge.QueryPointer()

					if err0 == nil {

						config.Lock()

						controller := config.GetMachine(config.Controller)

						if config.This == config.Controller {

							var target *types.Machine

							if mouse_x <= 1 {
								target = config.QueryMachine("left-of")
							} else if mouse_x >= controller.Screen.Width - 1 {
								target = config.QueryMachine("right-of")
							} else if mouse_y <= 1 {
								target = config.QueryMachine("above")
							} else if mouse_y >= controller.Screen.Height - 1 {
								target = config.QueryMachine("below")
							} else {
								target = config.QueryMachine("center")
							}

							if target != nil {

								// TODO: Send to Socket of target machine

							}

						}
						// Only evaluate boundary switching if no remote is active
						if global_state.Active == nil {

							if target != nil {
								global_state.Active = target
								fmt.Printf("Activated remote machine: %s (%s)\n", target.Hostname, target.Position)

								// Warp pointer slightly back inside host bounds
								if target.Position == "left-of" {
									bridge.WarpPointer(1, y)
								} else if target.Position == "right-of" {
									bridge.WarpPointer(hostWidth-2, y)
								} else if target.Position == "above" {
									bridge.WarpPointer(x, 1)
								} else if target.Position == "below" {
									bridge.WarpPointer(x, hostHeight-2)
								}
							}

						} else {

							// Forward the event to active machine via long-lived socket
							if global_state.Active.Socket != nil {
								evJSON, _ := json.Marshal(event)
								select {
								case global_state.Active.Socket <- evJSON:
								default:
									// channel full, drop event to avoid blocking
								}
							}

						}

						global_state.Unlock()

					}

					// Optional: always log locally
					data, _ := json.Marshal(event)
					fmt.Printf("Mouse: %+v\n", string(data))

				case event  := <-bridge.KeyboardEvents:


					global_state.Lock()

					if global_state.Active != nil && global_state.Active.Socket != nil {
						evJSON, _ := json.Marshal(event)
						select {
						case global_state.Active.Socket <- evJSON:
						default:
							// drop if channel full
						}
					}

					global_state.Unlock()

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

		bridge.Destroy()

		return nil

	} else if err0 != nil {
		return err0
	} else if err1 != nil {
		return err1
	} else {
		return nil
	}

}

