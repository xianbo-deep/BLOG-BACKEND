package bootstrap

import (
	"Blog-Backend/core"
	"Blog-Backend/internal/controller/admin"
)

func InitComponet() {
	// service初始化
	admin.InitAccessLogService(core.DB)
	admin.InitAnalysisService(core.DB)
	admin.InitCommentService()
	admin.InitDashboardService()
	admin.InitLoginService()
	admin.InitPerformanceService()
	admin.InitVisitorMapService()

	// TODO 封装dao和controller
}
