package types

type VirtualScreen struct {
	Width    uint               `json:"width"`
	Height   uint               `json:"height"`
	Machines map[string]*Screen `json:"machines"`
}

func NewVirtualScreen() *VirtualScreen {

	return &VirtualScreen{
		Width:    0,
		Height:   0,
		Machines: make(map[string]*Screen),
	}

}

func (vs *VirtualScreen) AddMachine(hostname string, screen *Screen) {

	if vs.Machines == nil {
		vs.Machines = make(map[string]*Screen)
	}

	vs.Machines[hostname] = screen

}

func (vs *VirtualScreen) GetMachine(hostname string) *Screen {

	if vs.Machines == nil {
		return nil
	}

	return vs.Machines[hostname]

}
