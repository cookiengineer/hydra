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

	if hostname == config.Controller && content_type == "application/json" && x_protocol == "hydra" {

		machine := config.GetMachine(hostname)

		if machine == nil {

			bytes, err0 := io.ReadAll(request.Body)

			if err0 == nil {

				var tmp types.Machine

				err1 := json.Unmarshal(bytes, &tmp)

				if err1 == nil {

					// Override IP with actual IP
					tmp.IP = request.RemoteAddr

					err2 := tmp.Parse()

					if err2 == nil {

						config.UpdateMachine(tmp)

						response.Header().Set("Content-Type",  "application/json")
						response.Header().Set("Cache-Control", "no-cache")
						response.Header().Set("Connection",    "keep-alive")

						flusher, ok := response.(http.Flusher)

						if ok == true {

							fmt.Println("Client connected: %s (%s)\n", tmp.Hostname, tmp.IP)

							for {
								select {
								case data := <-machine.Socket:
									fmt.Fprintf(response, "%s\n", data)
									flusher.Flush()
								case <-request.Context().Done():
									fmt.Println("Client disconnected: %s (%s)\n", tmp.Hostname, tmp.IP)
								case <-time.After(30 * time.Second):
									fmt.Fprintf(response, "{}\n")
									flusher.Flush()
								}
							}

						} else {

							response.WriteHeader(http.StatusUpgradeRequired)
							response.Write([]byte("{\"error\": \"Upgrade Required: Hydra Client must use keep-alive connections\"}"))

						}

					} else {

						response.Header().Set("Content-Type", "application/json")
						response.WriteHeader(http.StatusBadRequest)
						response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))

					}

				} else {

					response.Header().Set("Content-Type", "application/json")
					response.WriteHeader(http.StatusBadRequest)
					response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))

				}

			} else {

				response.Header().Set("Content-Type", "application/json")
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))

			}

		} else {

			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusConflict)
			response.Write([]byte("{\"error\": \"Conflict: Machine already registered\"}"))

		}

	} else {

		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusPreconditionFailed)
		response.Write([]byte("{\"error\": \"Precondition Failed: Not a Hydra Client\"}"))

	}

}
