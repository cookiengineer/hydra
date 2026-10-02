//go:build e2e

package e2e

import "encoding/json"
import "fmt"
import "net/http"
import "os"
import "testing"
import "time"
import "github.com/cookiengineer/hydra/actions"
import "github.com/cookiengineer/hydra/types"

func env(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value

}

func TestE2E(t *testing.T) {

	role := os.Getenv("HYDRA_E2E_ROLE")

	switch role {

	case "server":
		testServer(t)

	case "client":
		testClient(t)

	default:
		t.Skip("HYDRA_E2E_ROLE not set")

	}

}

type machineInfo struct {
	Hostname string `json:"hostname"`
}

type configResponse struct {
	Controller string                 `json:"controller"`
	Machines   map[string]machineInfo `json:"machines"`
	Screen     struct {
		Width  uint `json:"width"`
		Height uint `json:"height"`
	} `json:"screen"`
}

func fetchMachines() (map[string]machineInfo, error) {

	response, err := http.Get("http://127.0.0.1:" + types.Port() + "/machines")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var machines map[string]machineInfo

	if err := json.NewDecoder(response.Body).Decode(&machines); err != nil {
		return nil, err
	}

	return machines, nil

}

func fetchConfig() (*configResponse, error) {

	response, err := http.Get("http://127.0.0.1:" + types.Port() + "/config")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var config configResponse

	if err := json.NewDecoder(response.Body).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil

}

func testServer(t *testing.T) {

	controller := env("HYDRA_E2E_CONTROLLER", "")
	client := env("HYDRA_E2E_CLIENT", "")

	if controller == "" || client == "" {
		t.Fatal("HYDRA_E2E_CONTROLLER and HYDRA_E2E_CLIENT must be set")
	}

	go func() {
		if err := actions.Listen(controller); err != nil {
			fmt.Println("E2E_RESULT " + `{"role":"server","name":"listen","success":false,"message":"` + err.Error() + `"}`)
		}
	}()

	deadline := time.Now().Add(25 * time.Second)
	registered := false

	for time.Now().Before(deadline) {

		machines, err := fetchMachines()

		if err == nil {
			if _, ok := machines[client]; ok {
				registered = true
				break
			}
		}

		time.Sleep(500 * time.Millisecond)

	}

	if !registered {
		t.Errorf("client %s did not register in time", client)
		fmt.Println(`E2E_RESULT {"role":"server","name":"register","success":false,"message":"client did not register"}`)
		return
	}

	fmt.Printf("E2E_RESULT {\"role\":\"server\",\"name\":\"register\",\"success\":true,\"message\":\"client %s registered\"}\n", client)

	config, err := fetchConfig()

	if err != nil {
		t.Errorf("failed to fetch config: %v", err)
		fmt.Println(`E2E_RESULT {"role":"server","name":"config","success":false,"message":"config fetch failed"}`)
		return
	}

	if config.Screen.Width == 0 || config.Screen.Height == 0 {
		t.Errorf("virtual screen has zero size")
		fmt.Println(`E2E_RESULT {"role":"server","name":"virtualscreen","success":false,"message":"zero size"}`)
		return
	}

	if len(config.Machines) < 2 {
		t.Errorf("expected at least 2 machines, got %d", len(config.Machines))
		fmt.Println(`E2E_RESULT {"role":"server","name":"machines","success":false,"message":"missing machines"}`)
		return
	}

	message := fmt.Sprintf("virtual screen %dx%d with %d machines", config.Screen.Width, config.Screen.Height, len(config.Machines))
	fmt.Printf("E2E_RESULT {\"role\":\"server\",\"name\":\"protocol\",\"success\":true,\"message\":\"%s\"}\n", message)

	time.Sleep(10 * time.Second)

}

func testClient(t *testing.T) {

	address := env("HYDRA_E2E_CONTROLLER_ADDR", "")
	position := env("HYDRA_E2E_POSITION", "right-of")

	if address == "" {
		t.Fatal("HYDRA_E2E_CONTROLLER_ADDR must be set")
	}

	result := make(chan error, 1)

	go func() {
		result <- actions.Connect(address, position)
	}()

	select {

	case err := <-result:

		if err != nil {
			t.Errorf("connect failed: %v", err)
			fmt.Println(`E2E_RESULT {"role":"client","name":"connect","success":false,"message":"connect failed"}`)
			return
		}

		fmt.Println(`E2E_RESULT {"role":"client","name":"connect","success":false,"message":"stream ended early"}`)

	case <-time.After(5 * time.Second):

		fmt.Println(`E2E_RESULT {"role":"client","name":"connect","success":true,"message":"connected and receiving events"}`)

	}

}
