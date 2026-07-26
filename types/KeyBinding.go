package types

type KeyBinding struct {
	Modifiers uint32 `json:"modifiers"`
	Keycode   uint32 `json:"keycode"`
	Action    string `json:"action"`
}

const (
	ModShift   uint32 = 1 << 0
	ModLock    uint32 = 1 << 1
	ModControl uint32 = 1 << 2
	Mod1       uint32 = 1 << 3
	Mod2       uint32 = 1 << 4
	Mod3       uint32 = 1 << 5
	Mod4       uint32 = 1 << 6
	Mod5       uint32 = 1 << 7
)

const (
	SuperKeyMask uint32 = Mod4
)

const (
	XK_Escape       uint32 = 0xFF1B
	XK_Left         uint32 = 0xFF51
	XK_Up           uint32 = 0xFF52
	XK_Right        uint32 = 0xFF53
	XK_Down         uint32 = 0xFF54
	XK_Super_L      uint32 = 0xFFEB
	XK_Super_R      uint32 = 0xFFEC
	XK_Shift_L      uint32 = 0xFFE1
	XK_Shift_R      uint32 = 0xFFE2
	XK_Control_L    uint32 = 0xFFE3
	XK_Control_R    uint32 = 0xFFE4
	XK_Alt_L        uint32 = 0xFFE9
	XK_Alt_R        uint32 = 0xFFEA
)

const (
	ActionResetToController = "reset-to-controller"
	ActionTileLeft          = "tile-left"
	ActionTileRight         = "tile-right"
	ActionTileTopLeft       = "tile-top-left"
	ActionTileTopRight      = "tile-top-right"
	ActionTileBottomLeft    = "tile-bottom-left"
	ActionTileBottomRight   = "tile-bottom-right"
)

func GetDefaultKeyBindings() []KeyBinding {

	return []KeyBinding{
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Escape,
			Action:    ActionResetToController,
		},
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Left,
			Action:    ActionTileLeft,
		},
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Right,
			Action:    ActionTileRight,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Left,
			Action:    ActionTileTopLeft,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Up,
			Action:    ActionTileTopRight,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Down,
			Action:    ActionTileBottomLeft,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Right,
			Action:    ActionTileBottomRight,
		},
	}

}

func (kb *KeyBinding) Matches(modifiers uint32, keycode uint32) bool {
	return kb.Modifiers == modifiers && kb.Keycode == keycode
}
