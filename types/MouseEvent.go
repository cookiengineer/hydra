package types

type MouseEvent struct {
	Type   MouseEventType   `json:"type"`
	X      uint             `json:"x"`
	Y      uint             `json:"y"`
	DX     int              `json:"dx"`
	DY     int              `json:"dy"`
	Button MouseEventButton `json:"button"`
}
