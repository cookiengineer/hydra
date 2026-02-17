package types

type KeyboardEvent struct {
	Type    KeyboardEventType `json:"type"`
	Keycode uint32            `json:"keycode"`
}
