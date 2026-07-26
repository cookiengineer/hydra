package handlers

import "encoding/json"
import "fmt"
import "net/http"
import "github.com/cookiengineer/hydra/types"

func OnDisconnect(config *types.Config, state *types.GlobalState, response http.ResponseWriter, request *http.Request) {

	var machine types.Machine

	err := json.NewDecoder(request.Body).Decode(&machine)

	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("{\"error\": \"Bad Request: Invalid Payload\"}"))
		return
	}

	fmt.Printf("/disconnect from %s\n", machine.Hostname)

	config.RemoveMachine(machine)
	config.ComputeVirtualScreen()
	state.ResetActive()

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("{\"status\": \"ok\"}"))

}
