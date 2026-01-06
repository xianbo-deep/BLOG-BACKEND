package dao

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"context"
	"errors"
	"time"
)

// 在REDIS查今日最慢Top10
func GetSlowPages(ctx context.Context, limit int) ([]response.SlowDelayItem, error) {
	if core.RDB == nil {
		return nil, errors.New("Failed to get slow pages")
	}
	var result []response.SlowDelayItem
	// 获取键名
	today := time.Now().Format("2006-01-02")
	key := "blog:stat:daily:" + today + ":latency:rank"

	// 进行查询
	top10, err := core.RDB.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	// 处理返回结果
	for _, z := range top10 {
		result = append(result, response.SlowDelayItem{
			Path:     z.Member.(string),
			AvgDelay: int64(z.Score),
		})
	}

	return result, err
}

func GetAverageDelay() ([]response.AverageDelayItem, error) {
	// 将时间往前调整24h
	startTime := time.Now().Add(-24 * time.Hour)

	var res []response.AverageDelayItem

	db := core.DB.Model(&model.VisitLog{})

	err := db.Select("DATE_TRUNK('hour','visit_time') as hour , AVG(Latency) as avg_latency").
		Where("visit_time > ?", startTime).
		Group("DATE_TRUNK('hour','visit_time')").
		Order("hour ASC").
		Scan(&res).Error

	if err != nil {
		return res, err
	}

	return res, nil

}
