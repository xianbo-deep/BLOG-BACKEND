package deadlink

import (
	"Blog-Backend/bootstrap"
	"Blog-Backend/consts"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/model"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func RegisterDeadLink(c *cron.Cron, cmp *bootstrap.Components) {
	cfg := Config{
		SitemapURL:  os.Getenv(consts.EnvBaseURL) + sitemapURLSuffix,
		Retry:       retryTimes,
		Concurrency: defaultConcurrency,
		Timeout:     timeout,
	}

	mailer := cmp.Mailer

	checker := NewChecker(cfg, mailer)

	deadlinkDao := dao.NewDeadLinkDao(cmp.DB)

	// 注册定时任务
	_, err := c.AddFunc("0 0 0 * * *", func() {
		sum, res, err := checker.Check()
		if err != nil {
			log.Printf("[deadlink] err=%v", err)
			return
		}
		data := checker.processData(sum, res)

		// 加入数据库
		run := model.DeadLinkRun{
			StartedAt:    sum.StartedAT,
			FinishedAt:   sum.FinishedAT,
			PagesScanned: sum.PagesScanned,
			DeadLinkCnt:  sum.DeadlinkCnt,
			LinksChecked: sum.LinksChecked,
		}

		items := make([]model.DeadLinkItem, 0, len(res))
		for _, r := range res {
			items = append(items, model.DeadLinkItem{
				FromPage:   r.FromPage,
				LinkURL:    r.LinkURL,
				StatusCode: r.StatusCode,
				OK:         r.OK,
				Err:        r.Err,
				CheckedAt:  r.CheckedAT,
			})
		}
		if err := deadlinkDao.SaveRunAndItems(run, items); err != nil {
			log.Printf("死链检测数据插入数据库失败:%v", err)
		}

		// 发送邮箱通知
		err = mailer.SendTemplate([]string{consts.MyTencentEmail}, email.MailDeadlinkReport, data, true)
		if err != nil {
			log.Printf("[deadlink] err=%v", err)
		}

	})
	if err != nil {
		log.Printf("死链检测定时任务启动失败: %v", err)
	}
}
