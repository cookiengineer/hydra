package types

type Screen struct {
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Monitors []Monitor `json:"monitors"`
}
