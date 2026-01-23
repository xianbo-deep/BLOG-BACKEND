package bootstrap

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/internal/ws"
	"Blog-Backend/thirdparty/github"
	"Blog-Backend/thirdparty/github/service"
	"os"

	ctrl_admin "Blog-Backend/internal/controller/admin"
	ctrl_public "Blog-Backend/internal/controller/public"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/dao/cache"
	svc_admin "Blog-Backend/internal/service/admin"
	svc_public "Blog-Backend/internal/service/public"
)

type Components struct {
	Admin struct {
		AccessLog   *ctrl_admin.AccessLogController
		Analysis    *ctrl_admin.AnalysisController
		Comment     *ctrl_admin.CommentController
		Login       *ctrl_admin.LoginController
		Dashboard   *ctrl_admin.DashboardController
		Performance *ctrl_admin.PerformanceController
		VisitorMap  *ctrl_admin.VisitorMapController
		Websocket   *ctrl_admin.WebSocketController
	}
	Public struct {
		Collect *ctrl_public.CollectController
	}
}

func InitComponet() *Components {
	c := &Components{}
	// CacheClient初始化
	cacheClient := cache.NewCacheDAO(core.RDB)

	// GithubClient初始化
	client := github.NewClient(os.Getenv(consts.EnvDiscussionToken))

	// GithubDiscussionService初始化
	discussionService := service.NewDiscussionService(client)

	// 初始化websocket的hub并启动它
	hub := ws.NewHub()
	go hub.Run()

	// dao初始化
	analysisDao := dao.NewAnalysisDao(core.DB)
	collectDao := dao.NewCollectDao(core.DB, core.RDB)
	dashboardDao := dao.NewDashboardDao(core.DB, core.RDB)
	performanceDao := dao.NewPerformanceDao(core.DB, core.RDB)
	visitormapDao := dao.NewVisitorMapDao(core.DB)

	// service初始化
	accesslogService := svc_admin.NewAccessLogService(core.DB)
	analysisService := svc_admin.NewAnalysisService(analysisDao)
	commentService := svc_admin.NewCommentService(cacheClient, discussionService)
	dashboardService := svc_admin.NewDashboardService(dashboardDao)
	loginService := svc_admin.NewLoginService()
	performanceService := svc_admin.NewPerformanceService(performanceDao)
	visitormapService := svc_admin.NewVisitorMapSerive(visitormapDao)
	collectService := svc_public.NewCollectService(collectDao, hub)

	// controller初始化
	c.Admin.AccessLog = ctrl_admin.NewAccessLogController(accesslogService)
	c.Admin.Analysis = ctrl_admin.NewAnalysisController(analysisService)
	c.Admin.Comment = ctrl_admin.NewCommentController(commentService)
	c.Admin.Login = ctrl_admin.NewLoginController(loginService)
	c.Admin.Dashboard = ctrl_admin.NewDashboardController(dashboardService)
	c.Admin.Performance = ctrl_admin.NewPerformanceController(performanceService)
	c.Admin.VisitorMap = ctrl_admin.NewVisitorMapController(visitormapService)
	c.Admin.Websocket = ctrl_admin.NewWebSocketController(hub)
	c.Public.Collect = ctrl_public.NewCollectController(collectService)

	return c
}
