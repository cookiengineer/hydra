package actions

import "context"
import "encoding/json"
import "fmt"
import "net/http"
import "os"
import "os/signal"
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
		os.Interrupt,
		syscall.SIGTERM,
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
			IP:       "",
			Position: "center",
			Screen:   screen,
		})

		config.SetThis(host)
		state := types.NewGlobalState()
		state.SetScreen(config.Screen)

		http.HandleFunc("/config", func(response http.ResponseWriter, request *http.Request) {
			config.Mutex.Lock()
			defer config.Mutex.Unlock()
			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(config)
		})

		http.HandleFunc("/machines", func(response http.ResponseWriter, request *http.Request) {
			config.Mutex.Lock()
			defer config.Mutex.Unlock()
			response.Header().Set("Content-Type", "application/json")

			machines := make(map[string]interface{})
			for name, machine := range config.Machines {
				machines[name] = map[string]string{
					"hostname": machine.Hostname,
					"ip":       machine.IP,
					"position": machine.Position,
				}
			}
			json.NewEncoder(response).Encode(machines)
		})

		http.HandleFunc("/connect", func(response http.ResponseWriter, request *http.Request) {
			handlers.OnConnect(config, state, response, request)
		})

		http.HandleFunc("/disconnect", func(response http.ResponseWriter, request *http.Request) {
			handlers.OnDisconnect(config, state, response, request)
		})

		go http.ListenAndServe(":3000", nil)

		go bridge.Init()

		go func() {

			for {
				select {
				case event := <-bridge.MouseEvents:
					handleMouseEvent(bridge, event, state, config)

				case event := <-bridge.KeyboardEvents:
					handleKeyboardEvent(bridge, event, state, config)

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
