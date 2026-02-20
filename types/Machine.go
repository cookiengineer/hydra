package types

type Machine struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Position string `json:"position"` // left-of, right-of, above, below
	Screen   Screen `json:"screen"`
}
