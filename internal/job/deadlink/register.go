package deadlink

import (
	"Blog-Backend/consts"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func RegisterDeadLink(c *cron.Cron) {
	cfg := Config{
		SitemapURL:  os.Getenv(consts.EnvBaseURL) + sitemapURLSuffix,
		Retry:       retryTimes,
		Concurrency: defaultConcurrency,
		Timeout:     timeout,
	}
	checker := NewChecker(cfg)
	// 注册定时任务
	c.AddFunc("0 0 0 * * *", func() {
		sum, res, err := checker.Check()
		if err != nil {
			log.Printf("[deadlink] err=%v", err)
			return
		}
		// TODO 发送邮箱通知

		// TODO 加入数据库
	})
}
