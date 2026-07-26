package types

import "sync"

type GlobalState struct {
	Active            *Machine
	Mutex             sync.Mutex
	Screen            *VirtualScreen
	ActiveWorkspace   string
	Workspaces        map[string]*Workspace
	LastFocusedWindow uint64
}

func NewGlobalState() *GlobalState {

	workspaces := make(map[string]*Workspace)

	for _, ws := range GetDefaultWorkspaces() {
		ws.Windows = []Window{}
		workspaces[ws.Name] = &ws
	}

	return &GlobalState{
		Active:          nil,
		Mutex:           sync.Mutex{},
		Screen:          nil,
		ActiveWorkspace: "FG",
		Workspaces:      workspaces,
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

func (state *GlobalState) SetActiveWorkspace(name string) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.ActiveWorkspace = name

}

func (state *GlobalState) GetActiveWorkspace() string {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.ActiveWorkspace

}

func (state *GlobalState) StoreWorkspaceLayout(name string, windows []Window) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	if ws, ok := state.Workspaces[name]; ok {
		ws.Windows = windows
	}

}

func (state *GlobalState) GetWorkspaceLayout(name string) []Window {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	if ws, ok := state.Workspaces[name]; ok {
		return ws.Windows
	}

	return nil

}

func (state *GlobalState) GetWorkspace(name string) *Workspace {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	return state.Workspaces[name]

}

func (state *GlobalState) GetWorkspaceByName(name string) *Workspace {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	for _, ws := range state.Workspaces {
		if ws.Name == name {
			return ws
		}
	}

	return nil

}

func (state *GlobalState) GetWorkspaceByIndex(index uint32) *Workspace {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	for _, ws := range state.Workspaces {
		if ws.Index == index {
			return ws
		}
	}

	return nil

}

func (state *GlobalState) SetLastFocusedWindow(id uint64) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.LastFocusedWindow = id

}

func (state *GlobalState) GetLastFocusedWindow() uint64 {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.LastFocusedWindow

}
