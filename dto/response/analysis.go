package response

type AnalysisMetric struct {
	TotalPV    int `json:"totalPV"`
	TotalUV    int `json:"totalUV"`
	AvgLatency int `json:"avgLatency"`
	HotPage    int `json:"hotPage"`
	HotPagePV  int `json:"hotPagePV"`
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
	Referer string `json:"referers"`
	Percent int    `json:"percent"`
}

type AnalysisPathItemCountry struct {
	Country string `json:"country"`
	Percent int    `json:"percent"`
}

type AnalysisPathItemDetail struct {
	Path     string                    `json:"path"`
	Referers []AnalysisPathItemReferer `json:"referers"`
	Country  []AnalysisPathItemCountry `json:"country"`
}
