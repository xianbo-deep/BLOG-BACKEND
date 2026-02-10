package bootstrap

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/internal/ws"
	"Blog-Backend/thirdparty/github"
	"Blog-Backend/thirdparty/github/service"
	"os"

	ctrl_admin "Blog-Backend/internal/controller/admin"
	ctrl_github "Blog-Backend/internal/controller/github"
	ctrl_public "Blog-Backend/internal/controller/public"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/dao/cache"
	svc_admin "Blog-Backend/internal/service/admin"
	svc_github "Blog-Backend/internal/service/github"
	svc_public "Blog-Backend/internal/service/public"

	"gorm.io/gorm"
)

type Components struct {
	Mailer *email.Mailer

	GithubSVC *service.DiscussionService

	DB *gorm.DB

	Admin struct {
		AccessLog   *ctrl_admin.AccessLogController
		Analysis    *ctrl_admin.AnalysisController
		Comment     *ctrl_admin.CommentController
		Login       *ctrl_admin.LoginController
		Dashboard   *ctrl_admin.DashboardController
		Performance *ctrl_admin.PerformanceController
		VisitorMap  *ctrl_admin.VisitorMapController
		WebSocket   *ctrl_admin.WebSocketController
	}
	Public struct {
		Collect   *ctrl_public.CollectController
		Subscribe *ctrl_public.SubscribeController
	}
	Github struct {
		GithubWebhook *ctrl_github.GithubWebhookController
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

	c.GithubSVC = discussionService

	// 初始化websocket的hub并启动它
	hub := ws.NewHub()
	go hub.Run()

	// 初始化邮箱组件
	mailer := email.RegisterEmail()

	c.Mailer = mailer

	c.DB = core.DB

	// dao初始化
	analysisDao := dao.NewAnalysisDao(core.DB)
	collectDao := dao.NewCollectDao(core.DB, core.RDB)
	dashboardDao := dao.NewDashboardDao(core.DB, core.RDB)
	performanceDao := dao.NewPerformanceDao(core.DB, core.RDB)
	visitormapDao := dao.NewVisitorMapDao(core.DB)
	webhookDao := dao.NewGithubWebhookDao(core.DB)
	subscribeDao := dao.NewSubscribeDao(core.DB)

	// service初始化
	accesslogService := svc_admin.NewAccessLogService(core.DB)
	analysisService := svc_admin.NewAnalysisService(analysisDao)
	commentService := svc_admin.NewCommentService(cacheClient, discussionService)
	dashboardService := svc_admin.NewDashboardService(dashboardDao)
	loginService := svc_admin.NewLoginService()
	performanceService := svc_admin.NewPerformanceService(performanceDao)
	visitormapService := svc_admin.NewVisitorMapSerive(visitormapDao)
	collectService := svc_public.NewCollectService(collectDao, hub)
	githubWebhookService := svc_github.NewGithubWebhookService(discussionService, webhookDao)
	subscribeService := svc_public.NewSubscribeService(subscribeDao, mailer)

	// controller初始化
	c.Admin.AccessLog = ctrl_admin.NewAccessLogController(accesslogService)
	c.Admin.Analysis = ctrl_admin.NewAnalysisController(analysisService)
	c.Admin.Comment = ctrl_admin.NewCommentController(commentService)
	c.Admin.Login = ctrl_admin.NewLoginController(loginService)
	c.Admin.Dashboard = ctrl_admin.NewDashboardController(dashboardService)
	c.Admin.Performance = ctrl_admin.NewPerformanceController(performanceService)
	c.Admin.VisitorMap = ctrl_admin.NewVisitorMapController(visitormapService)
	c.Admin.WebSocket = ctrl_admin.NewWebSocketController(hub)
	c.Public.Collect = ctrl_public.NewCollectController(collectService)
	c.Github.GithubWebhook = ctrl_github.NewGithubWebhookController(githubWebhookService)
	c.Public.Subscribe = ctrl_public.NewSubscribeController(subscribeService)
	return c
}
