package middleware

import (
	"Blog-Backend/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")

		// token为空
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}

		// 格式校验
		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid pattern of token"})
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if claims.Username != os.Getenv("ADMIN_USER") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user,login again"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)

		c.Next()

	}
}
