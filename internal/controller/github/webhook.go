package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	g "Blog-Backend/dto/request/github"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/internal/service/github"
	"errors"
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
	err = ctrl.svc.NotifySubscribeUsers(pages, updatedAt, author)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, consts.CodeSuccess)
}

func latestCommitDocsPages(p g.PushPayload) (pages []email.ChangedPage, updatedAt time.Time, author string, err error) {
	// 获取提交作者
	author = p.Sender.Login

	headIdx := -1
	for i := range p.Commits {
		if p.Commits[i].ID == p.HeadCommit.ID {
			headIdx = i
			break
		}
	}
	if headIdx == -1 {
		if len(p.Commits) == 0 {
			return nil, consts.TransferTimeByLoc(time.Now()), author, errors.New("找不到最新的提交")
		}
		headIdx = len(p.Commits) - 1
	}

	haed := p.Commits[headIdx]

	// 解析时间
	ts := p.HeadCommit.Timestamp
	if ts == "" {
		ts = haed.Timestamp
	}
	if ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			updatedAt = t
		}
	}

	if updatedAt.IsZero() {
		updatedAt = consts.TransferTimeByLoc(time.Now())
	}

	// 获取变更文件
	base := os.Getenv(consts.EnvBaseURL)
	addPage := func(file, typ string) {
		if strings.HasPrefix(file, "docs/") && !strings.HasSuffix(file, "README.md") {
			path, ok := vuepressRouteHTML(file)
			if !ok {
				return
			}
			pages = append(pages, email.ChangedPage{
				Page:       file,
				ChangeType: typ,
				Path:       joinURL(base, path),
			})
		}
	}

	for _, f := range haed.Added {
		addPage(f, email.Added)
	}
	for _, f := range haed.Modified {
		addPage(f, email.Modified)
	}
	for _, f := range haed.Removed {
		addPage(f, email.Removed)
	}
	// 排序
	sort.Slice(pages, func(i, j int) bool { return pages[i].Path < pages[j].Path })
	return pages, updatedAt, author, nil
}

/* 工具函数 */

// 获取网页真实的相对路径
func vuepressRouteHTML(repoFile string) (string, bool) {
	if !strings.HasPrefix(repoFile, "docs/") {
		return "", false
	}
	rel := strings.TrimPrefix(repoFile, "docs/")

	// 根README.md
	if rel == "README.md" {
		return "/", true
	}

	// 目录README
	if strings.HasSuffix(rel, "README.md") {
		dir := strings.TrimSuffix(rel, "README.md")
		return "/" + dir, true
	}

	// 普通md
	if strings.HasSuffix(rel, ".md") {
		dir := strings.TrimSuffix(rel, ".md")
		return "/" + dir + ".html", true
	}
	return "", false
}

// 将域名和相对路径进行拼接
func joinURL(base, route string) string {
	base = strings.TrimRight(base, "/")
	if route == "" {
		return base + "/"
	}
	if !strings.HasPrefix(route, "/") {
		return base + "/" + route
	}
	return base + route
}
