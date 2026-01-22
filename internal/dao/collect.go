package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CollectDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewCollectDao(db *gorm.DB, rdb *redis.Client) *CollectDao {
	return &CollectDao{db: db, rdb: rdb}
}

func (d *CollectDao) InsertVisitLog(log model.VisitLog) error {
	/* 插入数据 */
	if d.db != nil {
		d.db.Create(&log)
		return nil
	}
	return errors.New("InsertVisitLog failed")
}

func (d *CollectDao) IncrementPV(ctx context.Context, path string) error {
	if d.rdb == nil {
		return errors.New("IncrementPV failed")
	}
	// 使用SortedSet维护排名
	pathPVKey := consts.GetDailyStatKey(time.Now().Format(consts.DateLayout), consts.RedisKeySuffixPathRank)
	totalPVKey := consts.GetDailyStatKey(time.Now().Format(consts.DateLayout), consts.RedisKeySuffixTotalPV)
	/* 插入数据 */
	if err := d.rdb.ZIncrBy(ctx, pathPVKey, 1, path).Err(); err != nil {
		return err
	}

	return d.rdb.IncrBy(ctx, totalPVKey, 1).Err()
}

func (d *CollectDao) IncrementUV(ctx context.Context, path string, visitorID string) error {
	if d.rdb == nil {
		return errors.New("IncrementUV failed")
	}
	// 键名
	// 单独页面 UV
	pathUVKey := consts.GetDailyPathUVKey(consts.GetTodayDate(), path)

	// 全站 UV
	totalUVKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixTotalUV)

	/* 插入数据 */
	if err := d.rdb.PFAdd(ctx, pathUVKey, visitorID).Err(); err != nil {
		return err
	}
	return d.rdb.PFAdd(ctx, totalUVKey, visitorID).Err()

}

func (d *CollectDao) RecordOnline(ctx context.Context, visitorID string) error {
	if d.rdb == nil {
		return errors.New("Failed to record online conut")
	}
	key := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixOnline)
	now := time.Now().Unix()

	// 添加新用户访问时间
	err := d.rdb.ZAdd(ctx, key, redis.Z{
		Member: visitorID,
		Score:  float64(now),
	}).Err()

	if err != nil {
		return err
	}
	// 删除过期用户
	cutoff := now - 3*60
	// 删除分数为0~cutoff的成员
	return d.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", cutoff)).Err()
}

// 插入访问页面延迟和页面访问次数
func (d *CollectDao) RecordLatency(ctx context.Context, path string, latency int64) error {
	if d.rdb == nil {
		return errors.New("Failed to record latency of path ")
	}
	countKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixPathCount)
	totalKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixPathTotalLatency)
	avgLatencyKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixPathAvgLatency)

	// 更新总延时
	if err := d.rdb.HIncrBy(ctx, totalKey, path, latency).Err(); err != nil {
		return err
	}

	// 更新访问次数
	if err := d.rdb.HIncrBy(ctx, countKey, path, 1).Err(); err != nil {
		return err
	}

	// 更新延时排名，方便取Top10
	total, _ := d.rdb.HGet(ctx, totalKey, path).Int64()
	count, _ := d.rdb.HGet(ctx, countKey, path).Int64()
	avg := float64(total) / float64(count)
	// 插入ZSet
	err := d.rdb.ZAdd(ctx, avgLatencyKey, redis.Z{
		Member: path,
		Score:  avg,
	}).Err()

	if err != nil {
		return err
	}

	return nil
}
