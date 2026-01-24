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

	ClientTime  time.Time
	IP          string
	Country     string
	CountryCode string
	CountryEN   string
	City        string
	CityEN      string
	Region      string
	UserAgent   string
	Referer     string
	RegionCode  string
	RegionEN    string
	Medium      string
	Source      string
	Device      string
	OS          string
	Browser     string
}
