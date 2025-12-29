package dao

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/response"
	"time"
)

// 获取过去六天的总访问量
func GetHistoryTrends(limit int) ([]response.DashboardTrends, error) {
	// 初始化变量
	var result []response.DashboardTrends
	// 给定日期
	today := time.Now().Format("2006-01-02")
	err := core.DB.Table("daily_article_stats").
		Select("to_char(date, 'YYYY-MM-DD') as date, sum(pv) as pv").
		Where("date < ?", today).
		Order("date desc").
		Limit(limit).
		Find(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

// 在Redis获取今天的访问量
func GetTodayPV() (response.DashboardTrends, error) {
	var result response.DashboardTrends
	today := time.Now().Format("2006-01-02")

	// TODO 去Redis查

	return result, nil
}

// 获取总日志数
func GetTotalLogs() (int64, error) {
	var result int64
	err := core.DB.Table("visit_logs").
		Count(&result).
		Error
	if err != nil {
		return 0, err
	}
	return result, nil
}

// 获取访问的来源
func GetGeoDistribution(
	startTime *time.Time,
	endTime *time.Time,
	limit *int,
) ([]response.GeoStatItem, error) {
	var result []response.GeoStatItem
	// 对limit进行赋值
	l := 5
	if limit != nil && *limit > 0 {
		l = *limit
	}

	// 初始化数据库语句
	db := core.DB.Table("visit_logs").
		Select("country,count(*) as count")

	if startTime != nil {
		db = db.Where("visit_time >= ?", *startTime)
	}

	if endTime != nil {
		db = db.Where("visit_time <= ?", *endTime)
	}

	err := db.
		Group("country").
		Order("count desc").
		Limit(l).
		Scan(&result).
		Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取错误日志
func GetErrorLogs(limit int) ([]response.ErrorLogItem, error) {
	var result []response.ErrorLogItem
	err := core.DB.Table("visit_logs").
		Select("path,status,visit_time").
		Where("status != 200").
		Order("visit_time desc").
		Limit(limit).
		Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
