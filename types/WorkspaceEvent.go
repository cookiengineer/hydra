package types

type WorkspaceEvent struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Index uint32 `json:"index"`
}
