package middleware

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/utils"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
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

func GithubWebhookVerify(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 raw body
		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// 放回去
		c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
		sig := c.GetHeader("X-Hub-Signature-256")
		if !verifyGitHubSignature(raw, sig, secret) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 把 raw body 放进 context，handler 里复用
		c.Set("raw_body", raw)

		c.Next()
	}
}

func verifyGitHubSignature(body []byte, sig string, secret string) bool {
	if secret == "" {
		return false
	}
	if !strings.HasPrefix(sig, "sha256=") {
		return false
	}
	// github提供的签名
	got := strings.TrimPrefix(sig, "sha256=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	// 自己算出来的签名
	expect := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expect), []byte(got))
}
