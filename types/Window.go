package types

type Window struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	MonitorID string `json:"monitor_id"`
	MachineID string `json:"machine_id"`
}
