package types

import "sync"

type GlobalState struct {
	Mutex    *sync.Mutex `json:"-"`
	Host     Machine     `json:"host"`
	Active   *Machine    `json:"active"`
	Machines []Machine   `json:"machines"`
	State    *State      `json:"state"`
}

func NewGlobalState(host Machine) *GlobalState {

	var state GlobalState

	state.Mutex    = &sync.Mutex{}
	state.Host     = host
	state.Machines = make([]Machine, 0)
	state.Machines = append(state.Machines, host)
	state.Active   = nil

	state.ComputeVirtualScreen()

	return &state

}

func (state *GlobalState) ComputeVirtualScreen() {

	virtual_screen := computeVirtualScreen(state.Host, state.Machines)

	if virtual_screen != nil {

		state.State = &State{
			VirtualScreen: virtual_screen,
		}

	} else {

		state.State = nil

	}

}
