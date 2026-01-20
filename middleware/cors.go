package middleware

import (
	"Blog-Backend/consts"
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	allow := map[string]bool{
		os.Getenv(consts.EnvBaseURL):  true,
		os.Getenv(consts.EnvAdminURL): true,
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allow[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin") // 重要：避免缓存把一个 origin 的结果给另一个
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
