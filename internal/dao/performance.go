package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PerformanceDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewPerformanceDao(db *gorm.DB, rdb *redis.Client) *PerformanceDao {
	return &PerformanceDao{db: db, rdb: rdb}
}

// 在REDIS查今日最慢Top10
func (d *PerformanceDao) GetSlowPages(ctx context.Context, limit int) ([]response.SlowDelayItem, error) {
	if d.rdb == nil {
		return nil, errors.New("Failed to get slow pages")
	}
	var result []response.SlowDelayItem
	// 获取键名
	today := time.Now().Format(consts.DateLayout)
	key := "blog:stat:daily:" + today + ":latency:rank"

	// 进行查询
	top10, err := d.rdb.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	// 处理返回结果
	for _, z := range top10 {
		var path string
		// 类型断言
		switch v := z.Member.(type) {
		case string:
			path = v
		case []byte:
			path = string(v)
		default:
			path = fmt.Sprint(v)
		}
		result = append(result, response.SlowDelayItem{
			Path:     path,
			AvgDelay: int64(z.Score),
		})
	}

	return result, err
}

func (d *PerformanceDao) GetAverageDelay() ([]response.AverageDelayItem, error) {
	// 将时间往前调整24h
	startTime := time.Now().Add(-consts.TimeRangeDay)

	var res []response.AverageDelayItem

	db := d.db.Model(&model.VisitLog{})

	err := db.Select("date_trunc('hour',visit_time) as time , avg(latency)::bigint as avg_delay").
		Where("visit_time > ?", startTime).
		Group("date_trunc('hour',visit_time)").
		Order("time ASC").
		Scan(&res).Error

	if err != nil {
		return res, err
	}

	for i := range res {
		res[i].Time = consts.TransferTimeByLoc(res[i].Time)
		res[i].Timestamp = consts.TransferTimeToTimestamp(res[i].Time)
	}
	return res, nil

}
