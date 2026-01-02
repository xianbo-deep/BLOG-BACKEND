package middleware

import (
	"Blog-Backend/consts"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 后续需要增加这个跨域域名配置
		c.Header("Access-Control-Allow-Origin", consts.EnvBaseURL)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, x-vercel-ip, x-vercel-ip-city")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
