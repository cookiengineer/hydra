package types

type MouseEvent struct {
	Type   MouseEventType   `json:"type"`
	DX     float64          `json:"dx"`
	DY     float64          `json:"dy"`
	Button MouseEventButton `json:"button"`
}
