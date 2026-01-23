package router

import (
	"Blog-Backend/bootstrap"
	"Blog-Backend/consts"

	"Blog-Backend/internal/controller/github"

	"Blog-Backend/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

/* 初始化路由器 */
func SetupRouter(c *bootstrap.Components) *gin.Engine {
	/* 设置为生产模式 */
	gin.SetMode(gin.ReleaseMode)

	/* 创建引擎 */
	r := gin.New()

	/* 配置受信任的代理 */
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil
	}

	/* 创建中间件 */
	r.Use(gin.Recovery())

	/* 使用中间件 */
	r.Use(
		middleware.CORSMiddleware(),
		middleware.TimeoutMiddleware(),
	)

	/* 定义路由组 */
	// 前端请求

	// 博客统计
	blogGroup := r.Group("/blog")
	// 使用中间件，简化业务
	blogGroup.Use(middleware.HeaderMiddleware())
	{
		// 统计流量
		blogGroup.Any("/collect", c.Public.Collect.CollectHandler)
	}

	// 后台统计
	adminGroup := r.Group("/admin")
	{
		// 登录
		adminGroup.POST("/login", c.Admin.Login.Login)

		// 鉴权
		adminAuth := adminGroup.Group("")
		adminAuth.Use(middleware.AuthMiddleware())
		{
			// 监控面板
			dashboard := adminAuth.Group("/dashboard")
			{
				dashboard.GET("/summary", c.Admin.Dashboard.GetDashboardSummary)
				dashboard.GET("/trend", c.Admin.Dashboard.GetDashboardTrend)
				dashboard.GET("/insights", c.Admin.Dashboard.GetDashboardInsights)
			}

			// 访问日志
			accesslog := adminAuth.Group("/accesslog")
			{
				accesslog.GET("/logs", c.Admin.AccessLog.GetAccessLogByQuery)
			}

			// 性能监控
			performance := adminAuth.Group("/performance")
			{
				performance.GET("/averageDelay", c.Admin.Performance.GetAverageDelay)
				performance.GET("/slowPages", c.Admin.Performance.GetSlowPages)
			}

			// 页面分析
			analysis := adminAuth.Group("/analysis")
			{
				analysis.GET("/metrics", c.Admin.Analysis.GetAnalysisMetrics)
				analysis.GET("/trend", c.Admin.Analysis.GetAnalysisTrend)
				analysis.GET("/rank", c.Admin.Analysis.GetAnalysisPathRank)
				analysis.GET("/path", c.Admin.Analysis.GetAnalysisPath)
				analysis.GET("/source", c.Admin.Analysis.GetAnalysisPathSource)
				analysis.GET("/querypath", c.Admin.Analysis.GetAnalysisPathByQuery)
				pathDetail := analysis.Group("/pathDetail")
				{
					pathDetail.GET("/trend", c.Admin.Analysis.GetAnalysisPathDetailTrend)
					pathDetail.GET("/metric", c.Admin.Analysis.GetAnalysisPathDetailMetric)
					pathDetail.GET("/source", c.Admin.Analysis.GetAnalysisPathDetailSource)
					pathDetail.GET("/device", c.Admin.Analysis.GetAnalysisPathDetailDevice)
				}
			}

			// 访客地图
			visitormap := adminAuth.Group("/visitormap")
			{
				visitormap.GET("/map", c.Admin.VisitorMap.GetVisitorMap)
				visitormap.GET("/chineseMap", c.Admin.VisitorMap.GetChineseVisitorMap)
			}

			// 评论区信息
			discussionmap := adminAuth.Group("/discussionmap")
			{
				discussionmap.GET("/metric", c.Admin.Comment.GetDiscussionMetric)
				discussionmap.GET("/trend", c.Admin.Comment.GetDiscussionTrend)
				discussionmap.GET("/activeuser", c.Admin.Comment.GetDiscussionActiveUser)
				discussionmap.GET("/feed", c.Admin.Comment.GetDiscussionNewFeed)
			}
		}
	}

	webhookGroup := r.Group("/webhook")
	{
		webhookGroup.POST("/github", middleware.GithubWebhookVerify(os.Getenv(consts.EnvGithubWebhookSecret)), github.GetNewNotify)
	}
	return r
}
