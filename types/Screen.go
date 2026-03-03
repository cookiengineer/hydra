package types

type Screen struct {
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Monitors []Monitor `json:"monitors"`
	OffsetX  int       `json:"offset_x"` // updated by computeVirtualScreen()
	OffsetY  int       `json:"offset_y"` // updated by computeVirtualScreen()
}
