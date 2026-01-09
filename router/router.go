package router

import (
	"Blog-Backend/internal/controller/admin"
	"Blog-Backend/internal/controller/public"
	"Blog-Backend/middleware"

	"github.com/gin-gonic/gin"
)

/* 初始化路由器 */
func SetupRouter() *gin.Engine {
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

	/* 配置跨域 */
	r.Use(middleware.CORSMiddleware())

	/* 定义路由组 */
	// 前端请求

	// 博客统计
	blogGroup := r.Group("/blog")
	{
		// 统计流量
		blogGroup.Any("/collect", public.CollectHandler)
	}

	// 后台统计
	adminGroup := r.Group("/admin")
	{
		// 登录
		adminGroup.POST("/login", admin.Login)

		// 鉴权
		adminAuth := adminGroup.Group("")
		adminAuth.Use(middleware.AuthMiddleware())
		{
			// 监控面板
			dashboard := adminAuth.Group("/dashboard")
			{
				dashboard.GET("/summary", admin.GetDashboardSummary)
				dashboard.GET("/trend", admin.GetDashboardTrend)
				dashboard.GET("/insights", admin.GetDashboardInsights)
			}

			// 访问日志
			accesslog := adminAuth.Group("/accesslog")
			{
				accesslog.GET("/logs", admin.GetAccessLog)
				accesslog.GET("/querylog", admin.GetAccessLogByQuery)
			}

			// 性能监控
			performance := adminAuth.Group("/performance")
			{
				performance.GET("/averageDelay", admin.GetAverageDelay)
				performance.GET("/slowPages", admin.GetSlowPages)
			}

			// 页面分析
			analysis := adminAuth.Group("/analysis")
			{
				analysis.GET("/metrics", admin.GetAnalysisMetrics)
				analysis.GET("/trend", admin.GetAnalysisTrend)
				analysis.GET("/rank", admin.GetAnalysisPathRank)
				analysis.GET("/path", admin.GetAnalysisPath)
				analysis.GET("/detail", admin.GetAnalysisPathDetail)
				analysis.GET("/querypath", admin.GetAnalysisPathByQuery)

			}

			// 访客地图
			visitormap := adminAuth.Group("/visitormap")
			{
				visitormap.GET("/map", admin.GetVisitorMap)
				visitormap.GET("/chineseMap", admin.GetChineseVisitorMap)
			}

			// 评论区信息
			discussionmap := adminAuth.Group("/discussionmap")
			{
				discussionmap.GET("/metric", admin.GetDiscussionMetric)
				discussionmap.GET("/trend", admin.GetDiscussionTrend)
				discussionmap.GET("/activeuser", admin.GetDiscussionActiveUser)
				discussionmap.GET("/feed", admin.GetDiscussionNewFeed)
			}
		}
	}

	return r
}
