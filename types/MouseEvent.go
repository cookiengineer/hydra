package types

type MouseEvent struct {
	Type   MouseEventType   `json:"type"`
	DX     int              `json:"dx"`
	DY     int              `json:"dy"`
	Button MouseEventButton `json:"button"`
}
