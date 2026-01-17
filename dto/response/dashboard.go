package response

type DashboardSummary struct {
	PV            int64   `json:"pv"`
	UV            int64   `json:"uv"`
	OnlineCount   int64   `json:"online_count"`
	TotalLogCount int64   `json:"total_log_count"`
	PvPercent     float64 `json:"pv_percent"`
	UVPercent     float64 `json:"uv_percent"`
}

type DashboardTrends struct {
	Date      string `json:"date"`
	PV        int64  `json:"pv"`
	UV        int64  `json:"uv"`
	Timestamp int64  `json:"timestamp"`
}

type GeoStatItem struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type ErrorLogItem struct {
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Time      string `json:"time" gorm:"column:visit_time"`
	Timestamp int64  `json:"timestamp"`
}

type DashboardInsights struct {
	GeoStats  []GeoStatItem  `json:"geo_stat_items"`
	ErrorLogs []ErrorLogItem `json:"error_logs"`
}
