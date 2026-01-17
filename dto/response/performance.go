package response

import "time"

type AverageDelayItem struct {
	Time      time.Time `json:"time"`
	AvgDelay  int64     `json:"avg_delay"`
	Timestamp int64     `json:"timestamp"`
}

type SlowDelayItem struct {
	Path     string `json:"path"`
	AvgDelay int64  `json:"avg_delay"`
}
