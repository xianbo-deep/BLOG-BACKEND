package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/github"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GithubWebhookController struct {
	svc *github.GithubWebhookService
}

func NewGithubWebhookController(svc *github.GithubWebhookService) *GithubWebhookController {
	return &GithubWebhookController{svc: svc}
}

// 有新评论
func (ctrl *GithubWebhookController) GetNewNotify(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")
	if event != "discussion" && event != "discussion_comment" {
		c.Status(200)
		return
	}

	if err := core.RDB.Incr(c, consts.RedisGithubCacheVerKey).Err(); err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}

	// 调用svc函数
	ctrl.svc.GetNewNotify(c.Request.Context())

	common.Success(c, consts.CodeSuccess)
}

// 通知订阅的用户
func (ctrl *GithubWebhookController) NotifySubscribeUser(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")

}
