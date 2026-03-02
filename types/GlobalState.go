package types

type GlobalState struct {
	Mutex    *sync.Mutex     `json:"-"`
	Host     Machine         `json:"host"`
	Active   *Machine        `json:"active"`
	Machines []types.Machine `json:"machines"`
	State    *State          `json:"state"`
}

func NewGlobalState(host Machine) *GlobalState {

	var state GlobalState

	state.Mutex    = &sync.Mutex{}
	state.Host     = host
	state.Machines = make([]types.Machine, 0)
	state.Machines = append(state.Machines, host)
	state.Active   = nil

	virtual_screen := math.ComputeVirtualScreen(state.Host, state.Machines)

	// TODO: Other State properties
	// TODO: Might have to reimplement listeners.Init into types.NewState()?

	state.State = &State{
		VirtualScreen: virtual_screen,
	}

	return &state

}
