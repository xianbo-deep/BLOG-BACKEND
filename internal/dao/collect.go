package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func InsertVisitLog(log model.VisitLog) error {
	/* 插入数据 */
	if core.DB != nil {
		core.DB.Create(&log)
		return nil
	}
	return errors.New("InsertVisitLog failed")
}

func IncrementPV(ctx context.Context, path string) error {
	if core.RDB == nil {
		return errors.New("IncrementPV failed")
	}
	// 使用SortedSet维护排名
	pathPVKey := consts.GetDailyStatKey(time.Now().Format(consts.DateLayout), consts.RedisKeySuffixPathRank)
	totalPVKey := consts.GetDailyStatKey(time.Now().Format(consts.DateLayout), consts.RedisKeySuffixTotalPV)
	/* 插入数据 */
	if err := core.RDB.ZIncrBy(ctx, pathPVKey, 1, path).Err(); err != nil {
		return err
	}

	return core.RDB.IncrBy(ctx, totalPVKey, 1).Err()
}

func IncrementUV(ctx context.Context, path string, visitorID string) error {
	if core.RDB == nil {
		return errors.New("IncrementUV failed")
	}
	// 键名
	// 单独页面 UV
	pathUVKey := consts.GetDailyPathUVKey(consts.GetTodayDate(), path)

	// 全站 UV
	totalUVKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixTotalUV)

	/* 插入数据 */
	if err := core.RDB.PFAdd(ctx, pathUVKey, visitorID).Err(); err != nil {
		return err
	}
	return core.RDB.PFAdd(ctx, totalUVKey, visitorID).Err()

}

func RecordOnline(ctx context.Context, visitorID string) error {
	if core.RDB == nil {
		return errors.New("Failed to record online conut")
	}
	key := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixOnline)
	now := time.Now().Unix()

	// 添加新用户访问时间
	err := core.RDB.ZAdd(ctx, key, redis.Z{
		Member: visitorID,
		Score:  float64(now),
	}).Err()

	if err != nil {
		return err
	}
	// 删除过期用户
	cutoff := now - 3*60
	// 删除分数为0~cutoff的成员
	return core.RDB.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", cutoff)).Err()
}

// 插入访问页面延迟和页面访问次数
func RecordLatency(ctx context.Context, path string, latency int64) error {
	if core.RDB == nil {
		return errors.New("Failed to record latency of path ")
	}
	countKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixPathCount)
	totalKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixPathTotalLatency)
	avgLatencyKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixPathAvgLatency)

	// 更新总延时
	if err := core.RDB.HIncrBy(ctx, totalKey, path, latency).Err(); err != nil {
		return err
	}

	// 更新访问次数
	if err := core.RDB.HIncrBy(ctx, countKey, path, 1).Err(); err != nil {
		return err
	}

	// 更新延时排名，方便取Top10
	total, _ := core.RDB.HGet(ctx, totalKey, path).Int64()
	count, _ := core.RDB.HGet(ctx, countKey, path).Int64()
	avg := float64(total) / float64(count)
	// 插入ZSet
	err := core.RDB.ZAdd(ctx, avgLatencyKey, redis.Z{
		Member: path,
		Score:  avg,
	}).Err()

	if err != nil {
		return err
	}

	return nil
}
