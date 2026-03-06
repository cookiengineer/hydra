package handlers

import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "strings"
import "time"
import "github.com/cookiengineer/hydra/types"

// OnConnect is the HTTP handler for new client connections.
// Keeps a long-lived line-based JSON socket in Machine.Socket.

func OnConnect(config *types.Config, response http.ResponseWriter, request *http.Request) {

	hostname     := strings.ToLower(strings.TrimSpace(request.Header.Get("Host")))
	content_type := strings.ToLower(strings.TrimSpace(request.Header.Get("Content-Type")))
	x_protocol   := strings.ToLower(strings.TrimSpace(request.Header.Get("X-Protocol")))
	ip           := request.RemoteAddr

	if hostname == config.Controller && content_type == "application/json" && x_protocol == "hydra" {

		machine := config.GetMachine(hostname)

		if machine == nil {

			bytes, err0 := io.ReadAll(request.Body)

			if err0 == nil {

				var tmp types.Machine

				err1 := json.Unmarshal(bytes, &tmp)

				if err1 == nil {

					err2 := tmp.Parse()

					if err2 == nil {

						config.UpdateMachine(tmp)

					} else {

						// TODO: Payload Error

					}

					config.UpdateMachine(types.Machine{
						Hostname: strings.ToLower(strings.TrimSpace(tmp.Hostname)),
						// TODO: Other properties
						// Socket: make(chan []byte),
					})
					// TODO: config.UpdateMachine(machine) ?

				} else {
					// TODO: Internal Server Error?
				}

			} else {
				// TODO: Internal Server Error?
			}

		} else {

			// TODO: Conflict?

		}

	} else {

		// TODO: Malformed payload

	}


		if machine == nil {
			global_state.Lock()
			global_state.Machines = append(global_state.Machines, types.Machine{
				Hostname: hostname,
				IP:       ip,
				Socket:   make(chan []byte, 128),
			})
			machine = &global_state.Machines[len(global_state.Machines)-1]
			global_state.Unlock()
		} else {
			// Reconnect: recreate socket channel
			machine.Socket = make(chan []byte, 128)
		}

		// Set headers for streaming
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Client connected: %s (%s)\n", hostname, ip)

		// Keep sending JSON lines until client disconnects
		for {
			select {
			case data := <-machine.Socket:
				// Send a line with JSON
				fmt.Fprintf(w, "%s\n", data)
				flusher.Flush()
			case <-r.Context().Done():
				fmt.Printf("Client disconnected: %s (%s)\n", hostname, ip)
				return
			case <-time.After(30 * time.Second):
				// Keep connection alive
				fmt.Fprintf(w, "{}\n")
				flusher.Flush()
			}
		}
	}
