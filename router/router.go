package router

import (
	"Blog-Backend/controller/admin"
	"Blog-Backend/controller/public"
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
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://xbzhong.cn")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, x-vercel-ip, x-vercel-ip-city")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	/* 定义路由组 */
	// 前端请求
	api := r.Group("/api")
	{
		// 统计流量
		api.Any("/collect", public.CollectHandler)

		// 后台统计
		adminGroup := api.Group("/admin")
		{
			// 监控面板
			dashboard := adminGroup.Group("/dashboard")
			{
				dashboard.GET("/summary", admin.GetDashboardSummary)

				dashboard.GET("/trend", admin.GetDashboardTrend)

				dashboard.GET("/insights", admin.GetDashboardInsights)
			}

			adminGroup.POST("/login", admin.Login)
		}

	}

	return r
}
