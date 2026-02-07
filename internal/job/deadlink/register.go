package deadlink

import (
	"Blog-Backend/consts"
	"Blog-Backend/internal/notify/email"
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

	mailer := email.RegisterEmail()

	checker := NewChecker(cfg, mailer)
	// 注册定时任务
	c.AddFunc("0 0 0 * * *", func() {
		sum, res, err := checker.Check()
		if err != nil {
			log.Printf("[deadlink] err=%v", err)
			return
		}
		data := checker.processData(sum, res)

		// TODO 加入数据库

		// TODO 发送邮箱通知
		err := mailer.SendTemplate([]string{}, email.MailDeadlinkReport, data)
		if err != nil {
			log.Printf("[deadlink] err=%v", err)
		}

	})
}
