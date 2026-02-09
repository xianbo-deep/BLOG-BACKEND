package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	g "Blog-Backend/dto/request/github"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/internal/service/github"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

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
	err := ctrl.svc.GetNewNotify(c.Request.Context())
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, consts.CodeSuccess)
}

// 通知订阅的用户
func (ctrl *GithubWebhookController) NotifySubscribeUsers(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")

	// 根据event选择不同的ctrl
	switch event {
	case "push":
		ctrl.handlePush(c)
	case "pull_request":
		ctrl.handlePullMerge(c)
	default:
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, errors.New("unknown event").Error())
		return
	}
}

// 处理push事件
func (ctrl *GithubWebhookController) handlePush(c *gin.Context) {
	var p g.PushPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}

	if p.Ref != consts.RepositoryRef {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, errors.New("Ref is invalid").Error())
		return
	}

	pages, updatedAt, author, err := latestCommitDocsPages(p)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	if len(pages) == 0 {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, "没有更改的文件")
		return
	}
	err = ctrl.svc.NotifySubscribeUser(pages, updatedAt, author)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, consts.CodeSuccess)
}

// 处理merge事件
func (ctrl *GithubWebhookController) handlePullMerge(c *gin.Context) {

}

func latestCommitDocsPages(p g.PushPayload) (pages []email.ChangedPage, updatedAt time.Time, author string, err error) {
	// 获取提交作者
	author = p.Sender.Login

	var head *struct {
		ID        string
		Timestamp string
		Added     []string
		Modified  []string
		Removed   []string
	}

	for i := range p.Commits {
		if p.Commits[i].ID == p.HeadCommit.ID {
			head = &p.Commits[i]
			break
		}
	}

	// 解析时间
	ts := p.HeadCommit.Timestamp
	if ts == "" && head != nil {
		ts = head.Timestamp
	}
	if ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			updatedAt = t
		}
	}

	if updatedAt.IsZero() {
		updatedAt = consts.TransferTimeByLoc(time.Now())
	}
	if head == nil {
		return nil, updatedAt, author, errors.New("找不到最新的提交")
	}

	// 获取变更文件
	base := os.Getenv(consts.EnvBaseURL)
	addPage := func(file, typ string) {
		if strings.HasPrefix(file, "docs/") && !strings.HasSuffix(file, "README.md") {
			path := strings.Trim(file, "docs")
			pages = append(pages, email.ChangedPage{
				Page:       file,
				ChangeType: typ,
				Path:       fmt.Sprintf("%s%s", base, path),
			})
		}
	}

	for _, f := range head.Added {
		addPage(f, email.Added)
	}
	for _, f := range head.Modified {
		addPage(f, email.Modified)
	}
	for _, f := range head.Removed {
		addPage(f, email.Removed)
	}
	// 排序
	sort.Slice(pages, func(i, j int) bool { return pages[i].Path < pages[j].Path })
	return pages, updatedAt, author, nil
}
