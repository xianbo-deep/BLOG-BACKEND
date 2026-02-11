package sync

import (
	"log"

	"github.com/robfig/cron/v3"
)

func RegisterSyncData(c *cron.Cron) {
	_, err := c.AddFunc("0 5 0 * * *", func() {
		log.Printf("注册数据同步定时任务")
		syncRedisToDB()
	})
	if err != nil {
		log.Printf("[Sync] err=%v", err)
	}
}
