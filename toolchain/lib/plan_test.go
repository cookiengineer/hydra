package lib

import "strings"
import "testing"

func testPlan() Plan {

	return Plan{
		Repo:            "/repo",
		Host:            "hydratwo",
		Image:           "hydra-e2e-builder",
		RemoteDir:       "/tmp/hydra-e2e",
		RemoteBin:       "/tmp/hydra-e2e/hydra-e2e.test",
		Controller:      "hydraone",
		Client:          "hydratwo",
		Address:         "192.168.0.11",
		Position:        "right-of",
		TestFilter:      "TestE2E",
		GuestDisplay:    ":0",
		GuestXauthority: "$HOME/.Xauthority",
	}

}

func TestBuildTestCommand(t *testing.T) {

	plan := testPlan()
	command := plan.BuildTestCommand()

	if command.Name != "podman" {
		t.Errorf("Expected podman, got %s", command.Name)
	}

	joined := strings.Join(command.Args, " ")

	if !strings.Contains(joined, "-tags=e2e") {
		t.Errorf("Expected e2e tag in %s", joined)
	}

	if !strings.Contains(joined, "./toolchain/e2e") {
		t.Errorf("Expected e2e package in %s", joined)
	}

	if !strings.Contains(joined, "/repo:/src") {
		t.Errorf("Expected repo mount in %s", joined)
	}

}

func TestRemoteClientCommand(t *testing.T) {

	plan := testPlan()
	command := plan.RemoteClientCommand()

	if command.Name != "ssh" {
		t.Errorf("Expected ssh, got %s", command.Name)
	}

	joined := strings.Join(command.Args, " ")

	if !strings.Contains(joined, "HYDRA_E2E_ROLE=client") {
		t.Errorf("Expected client role in %s", joined)
	}

	if !strings.Contains(joined, "HYDRA_E2E_CONTROLLER_ADDR=192.168.0.11") {
		t.Errorf("Expected controller address in %s", joined)
	}

	if !strings.Contains(joined, "DISPLAY=:0") {
		t.Errorf("Expected display in %s", joined)
	}

}

func TestLocalServerCommand(t *testing.T) {

	plan := testPlan()
	command := plan.LocalServerCommand()

	joined := strings.Join(command.Env, " ")

	if !strings.Contains(joined, "HYDRA_E2E_ROLE=server") {
		t.Errorf("Expected server role in %s", joined)
	}

	if command.Name != plan.LocalBinary() {
		t.Errorf("Expected local binary, got %s", command.Name)
	}

}
