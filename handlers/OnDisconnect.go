package handlers

import "encoding/json"
import "fmt"
import "net/http"
import "github.com/cookiengineer/hydra/types"

func OnDisconnect(config *types.Config, response http.ResponseWriter, request *http.Request) {

	var machine types.Machine

	json.NewDecoder(request.Body).Decode(&machine)

	fmt.Println("/disconnect from %s: %v", machine.Hostname, machine)

	config.UpdateMachine(machine)
	config.ComputeVirtualScreen()

	response.WriteHeader(http.StatusOK)

}
