package request

import "Blog-Backend/dto/common"

type AccessLogRequest struct {
	common.PageRequest
	IP        string `json:"ip" form:"ip"`
	Status    string `json:"status" form:"status"`
	VisitorID string `json:"visitor_id" form:"visitor_id"`
	Latency   int64  `json:"latency" form:"latency"`
	Path      string `json:"path" form:"path"`
}
