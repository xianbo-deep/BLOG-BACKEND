package request

type AnalysisRequest struct {
	Page     int    `form:"page" binding:"omitempty"`
	PageSize int    `form:"pageSize" binding:"omitempty"`
	Path     string `form:"path" binding:"omitempty"`
	Days     int    `form:"days" binding:"required"`
}
