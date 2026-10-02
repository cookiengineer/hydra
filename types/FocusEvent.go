package types

type FocusEvent struct {
	Type      string `json:"type"`
	Direction string `json:"direction,omitempty"`
	Edge      string `json:"edge,omitempty"`
}
