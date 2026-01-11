package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNewNotify(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")
	if event != "discussion" && event != "discussion_comment" {
		c.Status(200)
		return
	}

	if err := core.RDB.Incr(c, consts.RedisGithubCacheVerKey).Err(); err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}

	common.Success(c, consts.CodeSuccess)
}
