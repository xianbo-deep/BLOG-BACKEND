package request

// TODO 抽时间整理这些结构体标签的笔记
type AnalysisRequest struct {
	Page     int    `form:"page" binding:"omitempty"`
	PageSize int    `form:"pageSize" binding:"omitempty"`
	Path     string `form:"path" binding:"omitempty"`
	Days     int    `form:"days" binding:"required"`
}
