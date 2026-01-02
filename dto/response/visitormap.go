package response

type VisitorMapItem struct {
	Country  string `json:"country"`
	Visitors int64  `json:"visitors"`
}

type ChineseVisitorMapItem struct {
	Province string `json:"province"`
	Visitors int64  `json:"visitors"`
}
