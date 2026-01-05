package request

import "time"

type CollectRequest struct {
	VisitorID string `json:"visitor_id"` // 用户ID
	Path      string `json:"path"`       // 路径
	Status    int64  `json:"status"`     // 状态码
	Timestamp int64  `json:"timestamp"`  // 时间戳
	Latency   int64  `json:"latency"`    // 耗时
}

type CollectServiceDTO struct {
	VisitorID string
	Path      string
	Status    int64
	Latency   int64

	ClientTime time.Time
	IP         string
	Country    string
	City       string
	Region     string
	UserAgent  string
	Referer    string
}
