package handlers

import "net/http"
import "net/http/httptest"
import "strings"
import "testing"
import "github.com/cookiengineer/hydra/types"

func testConfigAndState() (*types.Config, *types.GlobalState) {

	config := types.NewConfig(types.Machine{
		Hostname: "controller",
		IP:       "10.0.0.1",
		Position: "center",
		Screen:   &types.Screen{Width: 1920, Height: 1080, Monitors: []types.Monitor{{Width: 1920, Height: 1080}}},
	})

	return config, types.NewGlobalState()

}

func TestOnConnectWrongProtocol(t *testing.T) {

	config, state := testConfigAndState()

	body := `{}`
	req := httptest.NewRequest("POST", "/connect", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Host", "controller")

	rr := httptest.NewRecorder()
	OnConnect(config, state, rr, req)

	if rr.Code != http.StatusPreconditionFailed {
		t.Errorf("Expected status 412, got %d", rr.Code)
	}

}

func TestOnConnectWrongHost(t *testing.T) {

	config, state := testConfigAndState()

	body := `{"hostname":"test","ip":"10.0.0.2"}`
	req := httptest.NewRequest("POST", "/connect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Protocol", "hydra")
	req.Header.Set("Host", "wrong-host")

	rr := httptest.NewRecorder()
	OnConnect(config, state, rr, req)

	if rr.Code != http.StatusPreconditionFailed {
		t.Errorf("Expected status 412, got %d", rr.Code)
	}

}

func TestOnConnectInvalidJSON(t *testing.T) {

	config, state := testConfigAndState()

	body := `not json`
	req := httptest.NewRequest("POST", "/connect", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Protocol", "hydra")
	req.Header.Set("Host", "controller")

	rr := httptest.NewRecorder()
	OnConnect(config, state, rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}

}

func TestOnDisconnect(t *testing.T) {

	config, state := testConfigAndState()

	config.UpdateMachine(types.Machine{
		Hostname: "disconnect-me",
		IP:       "10.0.0.5",
		Position: "left-of",
		Screen:   &types.Screen{Width: 1280, Height: 720, Monitors: []types.Monitor{{Width: 1280, Height: 720}}},
	})

	state.SetActive(config.GetMachine("disconnect-me"))

	body := `{"hostname":"disconnect-me","ip":"10.0.0.5","position":"left-of"}`
	req := httptest.NewRequest("POST", "/disconnect", strings.NewReader(body))
	rr := httptest.NewRecorder()

	OnDisconnect(config, state, rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	if config.GetMachine("disconnect-me") != nil {
		t.Error("Expected machine to be removed after disconnect")
	}

	if state.GetActive() != nil {
		t.Error("Expected active to be reset after disconnect")
	}

}

func TestOnDisconnectInvalidJSON(t *testing.T) {

	config, state := testConfigAndState()

	body := `not json`
	req := httptest.NewRequest("POST", "/disconnect", strings.NewReader(body))
	rr := httptest.NewRecorder()

	OnDisconnect(config, state, rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}

}
