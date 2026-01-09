package task

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/model"
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm/clause"
)

func InitCron() {
	// 创建cron
	c := cron.New(cron.WithSeconds())

	// 加入定时任务
	_, err := c.AddFunc("0 5 0 * * *", func() {
		log.Println("执行每日数据同步")
		SyncRedisToDB()
	})

	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
	}

	// 启动定时任务
	c.Start()

	log.Println("定时任务已启动")
}

func SyncRedisToDB() {
	ctx := context.Background()
	// 分布式锁
	lockKey := consts.RedisLockKey
	// 十分钟过期
	acquired, err := core.RDB.SetNX(ctx, lockKey, 1, 10*time.Minute).Result()
	if err != nil {
		log.Printf("获取Redis锁异常: %v", err)
		return
	}
	if !acquired {
		log.Printf("[Skip] 任务已被其他实例锁定")
		return
	}

	// 删除分布式锁
	defer core.RDB.Del(ctx, lockKey)

	yesterdayTime := time.Now().AddDate(0, 0, -1)
	yesterday := yesterdayTime.Format(consts.DateLayout)

	keyPathRank := consts.GetDailyStatKey(yesterday, consts.RedisKeySuffixPathRank)
	keyLatTotal := consts.GetDailyStatKey(yesterday, consts.RedisKeySuffixPathTotalLatency)
	keyLatCount := consts.GetDailyStatKey(yesterday, consts.RedisKeySuffixPathCount)
	keyTotalPV := consts.GetDailyStatKey(yesterday, consts.RedisKeySuffixTotalPV)
	keyTotalUV := consts.GetDailyStatKey(yesterday, consts.RedisKeySuffixTotalUV)
	keyAvgLatRank := consts.GetDailyStatKey(yesterday, consts.RedisKeySuffixPathAvgLatency)

	// 获取文章列表
	pathZSet, err := core.RDB.ZRevRangeWithScores(ctx, keyPathRank, 0, -1).Result()

	if err != nil {
		log.Printf("获取ZSet失败: %v", err)
		return
	}

	if len(pathZSet) == 0 {
		log.Printf("[Info] %s 无访问数据", yesterday)
		return
	}

	// 使用PipeLine实现批量读取
	pipe := core.RDB.Pipeline()

	// 创建切片存储返回结果
	uvCmds := make([]*redis.IntCmd, len(pathZSet))
	latTotalCmds := make([]*redis.StringCmd, len(pathZSet))
	latCountCmds := make([]*redis.StringCmd, len(pathZSet))

	for i, z := range pathZSet {
		path := z.Member.(string)
		uvKey := consts.GetDailyPathUVKey(yesterday, path)

		// 存储结果
		uvCmds[i] = pipe.PFCount(ctx, uvKey)
		latTotalCmds[i] = pipe.HGet(ctx, keyLatTotal, path)
		latCountCmds[i] = pipe.HGet(ctx, keyLatCount, path)
	}

	// 执行pipe
	_, err = pipe.Exec(ctx)

	if err != nil {
		log.Printf("无法从Redis中获取数据:%v", err)
		return
	}

	var stats []model.DailyArticleStat

	// 记录需要删除的key
	var uvKeysToDelete []string

	for i, z := range pathZSet {
		path := z.Member.(string)
		pv := int64(z.Score)

		// 解析并转换
		uv, _ := uvCmds[i].Result()
		totalLatStr, _ := latTotalCmds[i].Result()
		countLatStr, _ := latCountCmds[i].Result()

		totalLat, _ := strconv.ParseInt(totalLatStr, 10, 64)
		countLat, _ := strconv.ParseInt(countLatStr, 10, 64)

		stats = append(stats, model.DailyArticleStat{
			Path:         path,
			UV:           uv,
			PV:           pv,
			Date:         yesterdayTime,
			TotalLatency: totalLat,
			LatencyCount: countLat,
		})

		// 记录删除的key
		uvKeysToDelete = append(uvKeysToDelete, consts.GetDailyPathUVKey(yesterday, path))
	}

	// 批量插入数据库
	err := core.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}, {Name: "path"}},
		UpdateAll: true,
	}).Create(&stats).Error

	if err != nil {
		log.Printf("向数据库同步数据时发生错误:%v", err)
		return
	}

	/* 删除Redis的数据 */
	// 创建管道
	delPipe := core.RDB.Pipeline()

	// 删除文章维度的Key
	delPipe.Del(ctx, keyPathRank, keyLatTotal, keyLatCount)

	// 删除全站维度Key
	delPipe.Del(ctx, keyTotalPV, keyTotalUV, keyAvgLatRank)

	// 删除每篇文章独立的uvKey
	if len(uvKeysToDelete) > 0 {
		delPipe.Del(ctx, uvKeysToDelete...)
	}

	// 执行
	_, err = delPipe.Exec(ctx)
	if err != nil {
		log.Printf("删除redis中的key时出现错误:%v", err)
	} else {
		log.Printf("数据清理完成")
	}

}
