package task

import (
	"Blog-Backend/internal/job/sync"
	"log"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	// 创建cron
	c := cron.New(cron.WithSeconds())

	// 加入定时任务
	_, err := c.AddFunc("0 5 0 * * *", func() {
		log.Println("执行每日数据同步")
		sync.SyncRedisToDB()
	})

	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
	}

	// 启动定时任务
	c.Start()

	log.Println("定时任务已启动")
}
