package response

type AnalysisMetric struct {
	TotalPV    int64  `json:"totalPV"`
	TotalUV    int64  `json:"totalUV"`
	AvgLatency int64  `json:"avgLatency"`
	HotPage    string `json:"hotPage"`
	HotPagePV  int64  `json:"hotPagePV"`
}

type AnalysisTrendItem struct {
	Date string `json:"date"`
	PV   int64  `json:"pv"`
	UV   int64  `json:"uv"`
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

type AnalysisPathItemReferer struct {
	Referer string `json:"referer"`
	Percent int64  `json:"percent"`
}

type AnalysisPathItemCountry struct {
	Country string `json:"country"`
	Percent int64  `json:"percent"`
}

type AnalysisPathItemDetail struct {
	Path     string                    `json:"path"`
	Referers []AnalysisPathItemReferer `json:"referers"`
	Country  []AnalysisPathItemCountry `json:"country"`
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
	Date string `json:"date"`
	PV   int64  `json:"pv"`
	UV   int64  `json:"uv"`
}

type PathDetailSourceItem struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

type PathDetailDeviceItem struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}
