package types

import "sync"

type GlobalState struct {
	Mutex      *sync.Mutex         `json:"-"`
	Controller string              `json:"controller"` // has connected mouse and keyboard
	This       string              `json:"this"`       // this (local) machine
	Machines   map[string]*Machine `json:"machines"`   // all connected machines
	State      *State              `json:"state"`      // local state
}

func NewGlobalState(controller Machine) *GlobalState {

	var state GlobalState

	state.Mutex      = &sync.Mutex{}
	state.Controller = controller.Hostname
	state.This       = ""
	state.Machines   = make(map[string]*Machine)

	state.Machines[controller.Hostname] = &controller

	state.ComputeVirtualScreen()

	return &state

}

func (state *GlobalState) UpdateMachine(machine Machine) bool {

	if machine.Hostname != "" {

		state.Mutex.Lock()
		state.Machines[machine.Hostname] = &machine
		state.Mutex.Unlock()

		return true

	}

	return false

}

func (state *GlobalState) RemoveMachine(machine Machine) bool {

	if machine.Hostname != "" {

		_, ok := state.Machines[machine.Hostname]

		if ok == true {

			state.Mutex.Lock()
			delete(state.Machines, machine.Hostname)
			state.Mutex.Unlock()

		}

		return true

	}

	return false

}

func (state *GlobalState) SetThis(name string) bool {

	_, ok := state.Machines[name]

	if ok == true {

		state.Mutex.Lock()
		state.This = name
		state.Mutex.Unlock()

		state.ComputeVirtualScreen()

		return true

	}

	return false

}

func (state *GlobalState) ComputeVirtualScreen() {

	state.Mutex.Lock()

	virtual_screen := computeVirtualScreen(state.Controller, state.Machines)

	if virtual_screen != nil {

		state.State = &State{
			Screen: virtual_screen,
		}

	} else {

		state.State = nil

	}

	state.Mutex.Unlock()

}
