package types

import "sync"

type Config struct {
	Mutex      *sync.Mutex         `json:"-"`
	Controller string              `json:"controller"` // has connected mouse and keyboard
	This       string              `json:"this"`       // this (local) machine
	Machines   map[string]*Machine `json:"machines"`   // all connected machines
	Screen     *VirtualScreen      `json:"screen"`     // all connected screens
}

func NewConfig(controller Machine) *Config {

	var config Config

	config.Mutex      = &sync.Mutex{}
	config.Controller = controller.Hostname
	config.This       = ""
	config.Machines   = make(map[string]*Machine)

	config.Machines[controller.Hostname] = &controller

	config.ComputeVirtualScreen()

	return &config

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
