package types

type InitEvent struct {
	Type            string         `json:"type"`
	VirtualScreen   *VirtualScreen `json:"virtual_screen"`
	Workspaces      []Workspace    `json:"workspaces"`
	ActiveWorkspace string         `json:"active_workspace"`
}
