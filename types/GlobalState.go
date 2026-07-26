package types

import "sync"

type GlobalState struct {
	Active *Machine
	Mutex  sync.Mutex
	Screen *VirtualScreen
}

func NewGlobalState() *GlobalState {

	return &GlobalState{
		Active: nil,
		Mutex:  sync.Mutex{},
		Screen: nil,
	}

}

func (state *GlobalState) SetActive(machine *Machine) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.Active = machine

}

func (state *GlobalState) GetActive() *Machine {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.Active

}

func (state *GlobalState) ResetActive() {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.Active = nil

}

func (state *GlobalState) SetScreen(screen *VirtualScreen) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.Screen = screen

}

func (state *GlobalState) GetScreen() *VirtualScreen {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.Screen

}
