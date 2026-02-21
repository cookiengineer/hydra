package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cookiengineer/hydra/types"
)

// OnConnect is the HTTP handler for new client connections.
// Keeps a long-lived line-based JSON socket in Machine.Socket.
func OnConnect(global_state *types.GlobalState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		hostname := r.URL.Query().Get("hostname")
		if hostname == "" {
			http.Error(w, "missing hostname", http.StatusBadRequest)
			return
		}

		ip := r.RemoteAddr

		// Find existing machine or create new one
		var machine *types.Machine
		for i := range global_state.Machines {
			if global_state.Machines[i].Hostname == hostname {
				machine = &global_state.Machines[i]
				break
			}
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
}
