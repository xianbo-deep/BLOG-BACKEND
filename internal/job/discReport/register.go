package discReport

import (
	"Blog-Backend/bootstrap"
	"log"

	"github.com/robfig/cron/v3"
)

func RegisterDiscussionDigest(c *cron.Cron, cmp *bootstrap.Components) {
	discDigest := NewDiscussionDigest(cmp)
	_, err := c.AddFunc("0 0 8 * * MON", func() {
		log.Printf("注册评论区周报定时任务")
		discDigest.Start()

	})
	if err != nil {
		log.Printf("评论区周报定时任务启动失败: %v", err)
	}
}
