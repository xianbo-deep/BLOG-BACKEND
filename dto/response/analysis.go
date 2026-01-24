package response

import "time"

type AnalysisMetric struct {
	TotalPV    int64  `json:"totalPV"`
	TotalUV    int64  `json:"totalUV"`
	AvgLatency int64  `json:"avgLatency"`
	HotPage    string `json:"hotPage"`
	HotPagePV  int64  `json:"hotPagePV"`
}

type AnalysisTrendItem struct {
	Date      time.Time `json:"date"`
	PV        int64     `json:"pv"`
	UV        int64     `json:"uv"`
	Timestamp int64     `json:"timestamp"`
}

type AnalysisPathRankItem struct {
	Path string `json:"path"`
	PV   int64  `json:"pv"`
}

type AnalysisPathItem struct {
	Path       string `json:"path"`
	PV         int64  `json:"pv"`
	UV         int64  `json:"uv"`
	AvgLatency int    `json:"avgLatency"`
}

type AnalysisPathItemSource struct {
	Source  string `json:"source"`
	Percent int64  `json:"percent"`
}

type AnalysisPathItemDevice struct {
	Device  string `json:"device"`
	Percent int64  `json:"percent"`
}

type AnalysisPathItemDetail struct {
	Path     string                   `json:"path"`
	Referers []AnalysisPathItemSource `json:"sources"`
	Devices  []AnalysisPathItemDevice `json:"devices"`
}

type HotPageResult struct {
	Path string
	PV   int64
}

type PathDetailMetric struct {
	PV int64 `json:"pv"`
	UV int64 `json:"uv"`
}

type PathDetailTrendItem struct {
	Date      time.Time `json:"date"`
	PV        int64     `json:"pv"`
	UV        int64     `json:"uv"`
	Timestamp int64     `json:"timestamp"`
}

type PathDetailSourceItem struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

type PathDetailDeviceItem struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}
