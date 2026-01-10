package middleware

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
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
			common.Fail(c, http.StatusUnauthorized, consts.CodeTokenRequired, consts.ErrorMessage(consts.CodeTokenRequired))
			c.Abort()
			return
		}

		// 格式校验
		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.Fail(c, http.StatusUnauthorized, consts.CodeInvalidToken, consts.ErrorMessage(consts.CodeInvalidToken))
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])

		if err != nil {
			common.Fail(c, http.StatusUnauthorized, consts.CodeInvalidToken, err.Error())
			c.Abort()
			return
		}

		if claims.Username != os.Getenv("ADMIN_USER") {
			common.Fail(c, http.StatusUnauthorized, consts.CodeUserNotFound, consts.ErrorMessage(consts.CodeUserNotFound))
			c.Abort()
			return
		}

		c.Set("username", claims.Username)

		c.Next()

	}
}
