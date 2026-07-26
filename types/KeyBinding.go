package types

type KeyBinding struct {
	Modifiers uint32 `json:"modifiers"`
	Keycode   uint32 `json:"keycode"`
	Action    string `json:"action"`
	Data      uint32 `json:"data"`
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
	XK_grave        uint32 = 0x60
	XK_0            uint32 = 0x30
	XK_1            uint32 = 0x31
	XK_2            uint32 = 0x32
	XK_3            uint32 = 0x33
	XK_4            uint32 = 0x34
	XK_5            uint32 = 0x35
	XK_6            uint32 = 0x36
	XK_7            uint32 = 0x37
	XK_8            uint32 = 0x38
	XK_9            uint32 = 0x39
	XK_minus        uint32 = 0x2D
	XK_plus         uint32 = 0x2B
	XK_BackSpace    uint32 = 0xFF08
)

const (
	ActionResetToController = "reset-to-controller"
	ActionFocusLeft         = "focus-left"
	ActionFocusRight        = "focus-right"
	ActionFocusUp           = "focus-up"
	ActionFocusDown         = "focus-down"
	ActionTileLeft          = "tile-left"
	ActionTileRight         = "tile-right"
	ActionTileTop           = "tile-top"
	ActionTileBottom        = "tile-bottom"
	ActionTileTopLeft       = "tile-top-left"
	ActionTileTopRight      = "tile-top-right"
	ActionTileBottomLeft    = "tile-bottom-left"
	ActionTileBottomRight   = "tile-bottom-right"
	ActionSwitchWorkspace   = "switch-workspace"
	ActionMoveToWorkspace   = "move-to-workspace"
)

var workspaceKeys = []struct {
	Keycode uint32
	Index   uint32
}{
	{XK_grave, 0},
	{XK_0, 1},
	{XK_1, 2},
	{XK_2, 3},
	{XK_3, 4},
	{XK_4, 5},
	{XK_5, 6},
	{XK_6, 7},
	{XK_7, 8},
	{XK_8, 9},
	{XK_9, 10},
	{XK_minus, 11},
	{XK_plus, 12},
	{XK_BackSpace, 13},
}

func GetDefaultKeyBindings() []KeyBinding {

	bindings := []KeyBinding{
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Escape,
			Action:    ActionResetToController,
		},
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Left,
			Action:    ActionFocusLeft,
		},
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Right,
			Action:    ActionFocusRight,
		},
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Up,
			Action:    ActionFocusUp,
		},
		{
			Modifiers: SuperKeyMask,
			Keycode:   XK_Down,
			Action:    ActionFocusDown,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Left,
			Action:    ActionTileLeft,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Right,
			Action:    ActionTileRight,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Up,
			Action:    ActionTileTop,
		},
		{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   XK_Down,
			Action:    ActionTileBottom,
		},
	}

	for _, wk := range workspaceKeys {
		bindings = append(bindings, KeyBinding{
			Modifiers: SuperKeyMask,
			Keycode:   wk.Keycode,
			Action:    ActionSwitchWorkspace,
			Data:      wk.Index,
		})

		bindings = append(bindings, KeyBinding{
			Modifiers: SuperKeyMask | ModShift,
			Keycode:   wk.Keycode,
			Action:    ActionMoveToWorkspace,
			Data:      wk.Index,
		})
	}

	return bindings

}

func (kb *KeyBinding) Matches(modifiers uint32, keycode uint32) bool {
	return kb.Modifiers == modifiers && kb.Keycode == keycode
}
