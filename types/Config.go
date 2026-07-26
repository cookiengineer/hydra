package types

import "encoding/json"
import "os"
import "path/filepath"
import "sync"

type Config struct {
	Mutex       *sync.Mutex         `json:"-"`
	This        string              `json:"this"`
	Controller  string              `json:"controller"`
	Machines    map[string]*Machine `json:"machines"`
	Screen      *VirtualScreen      `json:"screen"`
	Workspaces  []Workspace         `json:"workspaces"`
	KeyBindings []KeyBinding        `json:"key_bindings"`
}

func GetDefaultWorkspaces() []Workspace {

	return []Workspace{
		{Name: "FG", Index: 0},
		{Name: "1", Index: 1},
		{Name: "2", Index: 2},
		{Name: "3", Index: 3},
		{Name: "4", Index: 4},
		{Name: "5", Index: 5},
		{Name: "6", Index: 6},
		{Name: "7", Index: 7},
		{Name: "8", Index: 8},
		{Name: "9", Index: 9},
		{Name: "10", Index: 10},
		{Name: "11", Index: 11},
		{Name: "12", Index: 12},
		{Name: "BG", Index: 13},
	}

}

func resolveConfigPath() string {

	config_home := os.Getenv("XDG_CONFIG_HOME")

	if config_home == "" {

		home, err := os.UserHomeDir()

		if err == nil {
			config_home = filepath.Join(home, ".config")
		}

	}

	if config_home == "" {
		config_home = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(config_home, "hydra", "config.json")

}

func LoadConfig() (*Config, error) {

	path := resolveConfigPath()

	data, err := os.ReadFile(path)

	var config Config

	if err == nil {
		err = json.Unmarshal(data, &config)
	}

	if err != nil || config.Mutex == nil {

		config = Config{
			Mutex:       &sync.Mutex{},
			Controller:  "",
			This:        "",
			Machines:    make(map[string]*Machine),
			Screen:      nil,
			Workspaces:  GetDefaultWorkspaces(),
			KeyBindings: GetDefaultKeyBindings(),
		}

		os.MkdirAll(filepath.Dir(path), 0755)
		json_data, _ := json.MarshalIndent(config, "", "  ")
		os.WriteFile(path, json_data, 0644)

	}

	return &config, nil

}

func NewConfig(controller Machine) *Config {

	config, _ := LoadConfig()

	config.Controller = controller.Hostname

	config.Machines[controller.Hostname] = &controller

	config.ComputeVirtualScreen()

	return config

}

func (config *Config) GetMachine(hostname string) *Machine {

	machine, ok := config.Machines[hostname]

	if ok == true {
		return machine
	}

	return nil

}

func (config *Config) QueryMachine(position string) *Machine {

	var result *Machine

	for _, machine := range config.Machines {

		if machine.Position == position {
			result = machine
			break
		}

	}

	return result

}

func (config *Config) UpdateMachine(machine Machine) bool {

	if machine.Hostname != "" {

		config.Mutex.Lock()
		config.Machines[machine.Hostname] = &machine
		config.Mutex.Unlock()

		return true

	}

	return false

}

func (config *Config) RemoveMachine(machine Machine) bool {

	if machine.Hostname != "" {

		_, ok := config.Machines[machine.Hostname]

		if ok == true {

			config.Mutex.Lock()
			delete(config.Machines, machine.Hostname)
			config.Mutex.Unlock()

		}

		return true

	}

	return false

}

func (config *Config) SetThis(name string) bool {

	_, ok := config.Machines[name]

	if ok == true {

		config.Mutex.Lock()
		config.This = name
		config.Mutex.Unlock()

		config.ComputeVirtualScreen()

		return true

	}

	return false

}

func (config *Config) ComputeVirtualScreen() {

	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	virtual_screen := computeVirtualScreen(config.Controller, config.Machines)

	if virtual_screen != nil {
		config.Screen = virtual_screen
	} else {
		config.Screen = nil
	}

}

func (config *Config) GetVirtualScreen() *VirtualScreen {

	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	return config.Screen

}
