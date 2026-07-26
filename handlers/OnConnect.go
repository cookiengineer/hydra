package handlers

import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "strings"
import "time"
import "github.com/cookiengineer/hydra/types"

func OnConnect(config *types.Config, state *types.GlobalState, response http.ResponseWriter, request *http.Request) {

	hostname     := strings.ToLower(strings.TrimSpace(request.Header.Get("Host")))
	content_type := strings.ToLower(strings.TrimSpace(request.Header.Get("Content-Type")))
	x_protocol   := strings.ToLower(strings.TrimSpace(request.Header.Get("X-Protocol")))

	if hostname != config.Controller || content_type != "application/json" || x_protocol != "hydra" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusPreconditionFailed)
		response.Write([]byte("{\"error\": \"Precondition Failed: Not a Hydra Client\"}"))
		return
	}

	bytes, err0 := io.ReadAll(request.Body)

	if err0 != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))
		return
	}

	var tmp types.Machine

	err1 := json.Unmarshal(bytes, &tmp)

	if err1 != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))
		return
	}

	tmp.IP = request.RemoteAddr

	err2 := tmp.Parse()

	if err2 != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))
		return
	}

	existing := config.GetMachine(tmp.Hostname)

	if existing != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusConflict)
		response.Write([]byte("{\"error\": \"Conflict: Machine already registered\"}"))
		return
	}

	config.UpdateMachine(tmp)
	config.ComputeVirtualScreen()

	machine := config.GetMachine(tmp.Hostname)

	if machine == nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("{\"error\": \"Internal Server Error\"}"))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")

	flusher, ok := response.(http.Flusher)

	if !ok {
		response.WriteHeader(http.StatusUpgradeRequired)
		response.Write([]byte("{\"error\": \"Upgrade Required: Hydra Client must use keep-alive connections\"}"))
		return
	}

	fmt.Printf("Client connected: %s (%s)\n", tmp.Hostname, tmp.IP)

	initPayload := map[string]interface{}{
		"type":           "init",
		"virtual_screen": config.GetVirtualScreen(),
	}
	initJSON, _ := json.Marshal(initPayload)
	fmt.Fprintf(response, "%s\n", initJSON)
	flusher.Flush()

	for {
		select {
		case data := <-machine.Socket:
			fmt.Fprintf(response, "%s\n", data)
			flusher.Flush()
		case <-request.Context().Done():
			fmt.Printf("Client disconnected: %s (%s)\n", tmp.Hostname, tmp.IP)
			config.RemoveMachine(tmp)
			config.ComputeVirtualScreen()
			state.ResetActive()
			return
		case <-time.After(30 * time.Second):
			fmt.Fprintf(response, "{}\n")
			flusher.Flush()
		}
	}

}
