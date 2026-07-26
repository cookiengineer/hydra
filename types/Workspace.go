package types

type Workspace struct {
	Name    string   `json:"name"`
	Index   uint32   `json:"index"`
	Windows []Window `json:"windows"`
}
