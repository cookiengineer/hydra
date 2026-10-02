package types

import "sync"

type GlobalState struct {
	Active            *Machine
	Mutex             sync.Mutex
	Screen            *VirtualScreen
	ActiveWorkspace   string
	Workspaces        map[string]*Workspace
	LastFocusedWindow uint64
	CursorX           int
	CursorY           int
	PressedKeys       map[uint32]bool
	PressedButtons    map[int]bool
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
		PressedKeys:     make(map[uint32]bool),
		PressedButtons:  make(map[int]bool),
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

func (state *GlobalState) SetCursor(x int, y int) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.CursorX = x
	state.CursorY = y

}

func (state *GlobalState) GetCursor() (int, int) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.CursorX, state.CursorY

}

func (state *GlobalState) TrackKey(keysym uint32) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	if state.PressedKeys == nil {
		state.PressedKeys = make(map[uint32]bool)
	}

	state.PressedKeys[keysym] = true

}

func (state *GlobalState) UntrackKey(keysym uint32) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	delete(state.PressedKeys, keysym)

}

func (state *GlobalState) IsTrackedKey(keysym uint32) bool {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.PressedKeys[keysym]

}

func (state *GlobalState) TrackedKeys() []uint32 {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	keys := make([]uint32, 0, len(state.PressedKeys))

	for keysym := range state.PressedKeys {
		keys = append(keys, keysym)
	}

	return keys

}

func (state *GlobalState) ClearTrackedKeys() {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.PressedKeys = make(map[uint32]bool)

}

func (state *GlobalState) TrackButton(button int) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	if state.PressedButtons == nil {
		state.PressedButtons = make(map[int]bool)
	}

	state.PressedButtons[button] = true

}

func (state *GlobalState) UntrackButton(button int) {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	delete(state.PressedButtons, button)

}

func (state *GlobalState) IsTrackedButton(button int) bool {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	return state.PressedButtons[button]

}

func (state *GlobalState) TrackedButtons() []int {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	buttons := make([]int, 0, len(state.PressedButtons))

	for button := range state.PressedButtons {
		buttons = append(buttons, button)
	}

	return buttons

}

func (state *GlobalState) ClearTrackedButtons() {

	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	state.PressedButtons = make(map[int]bool)

}
