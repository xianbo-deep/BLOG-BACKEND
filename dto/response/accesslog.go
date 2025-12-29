package response

import "time"

type AccessLog struct {
	Path       string    `json:"path"`
	VisitTime  time.Time `json:"visit_time"`
	IP         string    `json:"ip"`
	ClientTime time.Time `json:"client_time"`
	UserAgent  string    `json:"user_agent"`
	Referer    string    `json:"referer"`
	Country    string    `json:"country"`
	City       string    `json:"city"`
	Region     string    `json:"region"`
	Status     int64     `json:"status"`
}
