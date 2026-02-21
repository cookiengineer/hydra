package actions

import "context"
import "encoding/json"
import "fmt"
import "net/http"
import "os"
import "os/signal"
import "sync"
import "syscall"
import "github.com/cookiengineer/hydra/handlers"
import "github.com/cookiengineer/hydra/listeners"
import "github.com/cookiengineer/hydra/math"
import "github.com/cookiengineer/hydra/parsers"
import "github.com/cookiengineer/hydra/types"

type GlobalState struct {
	sync.Mutex
	Host          types.Machine   `json:"host"`
	Active        *types.Machine  `json:"active"`
	Machines      []types.Machine `json:"machines"`
	VirtualScreen types.Screen    `json:"virtual_screen"`
}

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
	screen, err1 := parsers.Xrandr()

	if err0 == nil && err1 == nil {

		host_machine := types.Machine{
			Hostname: host,
			IP:       "", // populated later
			Position: "host",
			Screen:   *screen,
		}

		global_state := &GlobalState{
			Host:     host_machine,
			Machines: make([]types.Machine, 0),
			Active:   nil,
		}

		global_state.VirtualScreen = math.ComputeVirtualScreen(global_state.Host, global_state.Machines)

		http.HandleFunc("/state", func(response http.ResponseWriter, request *http.Request) {

			global_state.Lock()

			// TODO: Error handling

			json.NewEncoder(response).Encode(global_state)

			global_state.Unlock()

		})

		http.HandleFunc("/connect", handlers.OnConnect(global_state))

		http.HandleFunc("/disconnect", func(response http.ResponseWriter, request *http.Request) {

			var machine types.Machine

			json.NewDecoder(request.Body).Decode(&machine)

			fmt.Println("/disconnect from %s: %v", machine.Hostname, machine)

			global_state.Lock()

			for m, other := range global_state.Machines {

				if other.Hostname == machine.Hostname {
					global_state.Machines = append(global_state.Machines[:m], global_state.Machines[m+1:]...)
					break
				}

			}

			global_state.VirtualScreen = math.ComputeVirtualScreen(global_state.Host, global_state.Machines)

			response.WriteHeader(http.StatusOK)

			global_state.Unlock()

		})

		go http.ListenAndServe(":3000", nil)

		go listeners.StartLoop(state)

		go func() {

			for {
				select {
				case <-state.MouseEvents:

					x, y, err := state.QueryPointer()

					if err == nil {

						global_state.Lock()

						hostWidth := global_state.Host.Screen.Width
						hostHeight := global_state.Host.Screen.Height

						// Only evaluate boundary switching if no remote is active
						if global_state.Active == nil {

							var target *types.Machine

							if x <= 0 {
								for i := range global_state.Machines {
									if global_state.Machines[i].Position == "left-of" {
										target = &global_state.Machines[i]
										break
									}
								}
							} else if x >= hostWidth-1 {
								for i := range global_state.Machines {
									if global_state.Machines[i].Position == "right-of" {
										target = &global_state.Machines[i]
										break
									}
								}
							} else if y <= 0 {
								for i := range global_state.Machines {
									if global_state.Machines[i].Position == "above" {
										target = &global_state.Machines[i]
										break
									}
								}
							} else if y >= hostHeight-1 {
								for i := range global_state.Machines {
									if global_state.Machines[i].Position == "below" {
										target = &global_state.Machines[i]
										break
									}
								}
							}

							if target != nil {
								global_state.Active = target
								fmt.Printf("Activated remote machine: %s (%s)\n", target.Hostname, target.Position)

								// Warp pointer slightly back inside host bounds
								if target.Position == "left-of" {
									state.WarpPointer(1, y)
								} else if target.Position == "right-of" {
									state.WarpPointer(hostWidth-2, y)
								} else if target.Position == "above" {
									state.WarpPointer(x, 1)
								} else if target.Position == "below" {
									state.WarpPointer(x, hostHeight-2)
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

				case event  := <-state.KeyboardEvents:


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

		state.Destroy()

		return nil

	} else if err0 != nil {
		return err0
	} else if err1 != nil {
		return err1
	} else {
		return nil
	}

}

