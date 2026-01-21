package bootstrap

import (
	"Blog-Backend/core"
	"Blog-Backend/internal/controller/admin"
)

// TODO 对GithubService也在这做初始化 还有一个Client
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
