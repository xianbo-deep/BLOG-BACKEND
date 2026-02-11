package task

import (
	"Blog-Backend/bootstrap"
	"Blog-Backend/internal/job/deadlink"
	"Blog-Backend/internal/job/discReport"
	"Blog-Backend/internal/job/sync"
	"log"

	"github.com/robfig/cron/v3"
)

func InitCron(cmp *bootstrap.Components) {
	// 创建cron
	c := cron.New(cron.WithSeconds())

	// 注册数据同步
	sync.RegisterSyncData(c)

	// 注册死链检测
	deadlink.RegisterDeadLink(c, cmp)

	// 注册评论区周报
	discReport.RegisterDiscussionDigest(c, cmp)

	// 启动定时任务
	c.Start()

	log.Println("定时任务已启动")
}
