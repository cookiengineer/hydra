package types

type HandoffEvent struct {
	Type      string `json:"type"`
	Machine   string `json:"machine"`
	Direction string `json:"direction"`
}
