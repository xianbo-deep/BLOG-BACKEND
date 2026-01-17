package response

type AverageDelayItem struct {
	Time      string `json:"time"`
	AvgDelay  int64  `json:"avg_delay"`
	Timestamp int64  `json:"timestamp"`
}

type SlowDelayItem struct {
	Path     string `json:"path"`
	AvgDelay int64  `json:"avg_delay"`
}
