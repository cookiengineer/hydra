package types

type Monitor struct {
	Output     string        `json:"output"`
	Connected  bool          `json:"connected"`
	Resolution string        `json:"resolution"`
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	OffsetX    int           `json:"offset_x"`
	OffsetY    int           `json:"offset_y"`
}
