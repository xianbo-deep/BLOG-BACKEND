package request

import "Blog-Backend/dto/common"

type AccessLogRequest struct {
	common.PageRequest
	KeyWord string `json:"keyword,omitempty" form:"keyword"`
	Status  string `json:"status,omitempty" form:"status"`
	Latency int64  `json:"latency" form:"latency"`
}
