package types

type Screen struct {
	Width    uint      `json:"width"`
	Height   uint      `json:"height"`
	Monitors []Monitor `json:"monitors"`
	OffsetX  uint      `json:"offset_x"` // updated by computeVirtualScreen()
	OffsetY  uint      `json:"offset_y"` // updated by computeVirtualScreen()
}
