package handlers

import "encoding/json"
import "fmt"
import "io"
import "net"
import "net/http"
import "strings"
import "time"
import "github.com/cookiengineer/hydra/types"

func normalizeHost(host string) string {

	host = strings.ToLower(strings.TrimSpace(host))

	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host

}

func OnConnect(config *types.Config, state *types.GlobalState, response http.ResponseWriter, request *http.Request) {

	host_header        := normalizeHost(request.Header.Get("Host"))
	controller_header  := normalizeHost(request.Header.Get("X-Hydra-Controller"))
	content_type       := strings.ToLower(strings.TrimSpace(request.Header.Get("Content-Type")))
	x_protocol         := strings.ToLower(strings.TrimSpace(request.Header.Get("X-Protocol")))

	requester := controller_header

	if requester == "" {
		requester = host_header
	}

	if requester != config.Controller || !strings.HasPrefix(content_type, "application/json") || x_protocol != "hydra" {
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

	if remote_host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil && remote_host != "" {
		tmp.IP = remote_host
	}

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

	init_event := types.InitEvent{
		Type:            "init",
		VirtualScreen:   config.GetVirtualScreen(),
		Workspaces:      config.Workspaces,
		ActiveWorkspace: state.GetActiveWorkspace(),
	}

	envelope, err3 := types.NewEvent("init", init_event)

	if err3 == nil {

		init_payload, _ := json.Marshal(envelope)
		fmt.Fprintf(response, "%s\n", init_payload)
		flusher.Flush()

	}

	for {
		select {
		case data := <-machine.Socket:
			fmt.Fprintf(response, "%s\n", data)
			flusher.Flush()
		case <-request.Context().Done():
			fmt.Printf("Client disconnected: %s (%s)\n", tmp.Hostname, tmp.IP)
			config.RemoveMachine(tmp)
			config.ComputeVirtualScreen()
			state.ClearTrackedKeys()
			state.ClearTrackedButtons()
			state.ResetActive()
			return
		case <-time.After(30 * time.Second):
			fmt.Fprintf(response, "{}\n")
			flusher.Flush()
		}
	}

}
