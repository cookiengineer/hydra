package types

type KeyboardEvent struct {
	Type      KeyboardEventType `json:"type"`
	Keycode   uint32            `json:"keycode,omitempty"`
	Keysym    uint32            `json:"keysym,omitempty"`
	Modifiers uint32            `json:"modifiers,omitempty"`
}
