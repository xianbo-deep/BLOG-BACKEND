package router

import (
	admin2 "Blog-Backend/internal/controller/admin"
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

	/* 创建中间件 */
	r.Use(gin.Recovery())

	/* 配置跨域 */
	r.Use(middleware.CORSMiddleware())

	/* 定义路由组 */
	// 前端请求
	api := r.Group("/api")
	{
		// 博客统计
		blogGroup := api.Group("/blog")
		{
			// 统计流量
			blogGroup.Any("/collect", public.CollectHandler)
		}

		// 后台统计
		adminGroup := api.Group("/admin")
		{
			// 登录
			adminGroup.POST("/login", admin2.Login)

			// 鉴权
			adminAuth := adminGroup.Group("")
			adminAuth.Use(middleware.AuthMiddleware())
			{
				// 监控面板
				dashboard := adminAuth.Group("/dashboard")
				{
					dashboard.GET("/summary", admin2.GetDashboardSummary)
					dashboard.GET("/trend", admin2.GetDashboardTrend)
					dashboard.GET("/insights", admin2.GetDashboardInsights)
				}

				// 访问日志
				accesslog := adminAuth.Group("/accesslog")
				{
					accesslog.GET("/logs", admin2.GetAccessLog)
				}

				// 性能监控
				performance := adminAuth.Group("/performance")
				{
					performance.GET("/averageDelay", admin2.GetAverageDelay)
					performance.GET("/slowPages", admin2.GetSlowPages)
				}

				// 页面分析
				analysis := adminAuth.Group("/analysis")
				{
					analysis.GET("/total", admin2.GetTotalPagesData)
					analysis.GET("/today", admin2.GetTodayPagesData)
				}
			}
		}

	}

	return r
}
