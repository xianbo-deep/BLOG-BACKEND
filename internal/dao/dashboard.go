package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DashboardDao struct {
	rdb *redis.Client
	db  *gorm.DB
}

func NewDashboardDao(db *gorm.DB, rdb *redis.Client) *DashboardDao {
	return &DashboardDao{db: db, rdb: rdb}
}

// 获取实时在线人数
func (d *DashboardDao) GetOnlineCount(ctx context.Context) (int64, error) {
	if d.rdb == nil {
		return 0, errors.New("Redis not initialized")
	}

	key := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixOnline)
	// 删除过期用户
	now := time.Now().Unix()
	start := now - 3*60
	// 删除分数为0~cutoff的成员
	if err := d.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", start)).Err(); err != nil {
		return 0, err
	}
	// 查询在线人数
	return d.rdb.ZCount(ctx, key, fmt.Sprintf("%d", start), fmt.Sprintf("%d", now)).Result()
}

// 获取今日PV和UV
func (d *DashboardDao) GetTodayPVUV(ctx context.Context) (int64, int64, error) {
	if d.rdb == nil {
		return 0, 0, errors.New("Failed to get pv and uv")
	}
	PVKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixTotalPV)
	UVKey := consts.GetDailyStatKey(consts.GetTodayDate(), consts.RedisKeySuffixTotalUV)

	pvStr, err := d.rdb.Get(ctx, PVKey).Result()

	if err != nil {
		if err == redis.Nil {
			pvStr = "0"
		} else {
			return 0, 0, err
		}
	}

	// 字符串转整型
	pv, _ := strconv.ParseInt(pvStr, 10, 64)

	uv, err := d.rdb.PFCount(ctx, UVKey).Result()
	if err != nil {
		return 0, 0, err
	}
	return uv, pv, nil
}

// 获取过去六天的总访问量
func (d *DashboardDao) GetHistoryTrends(limit int) ([]response.DashboardTrends, error) {
	// 初始化变量
	var result []response.DashboardTrends
	// 给定日期
	today := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay)
	err := d.db.Table("daily_article_stats").
		Select("to_char(date at time zone 'Asia/Shanghai', 'YYYY-MM-DD') as date, sum(pv) as pv,sum(uv) as uv").
		Where("date < ?", today).
		Order("to_char(date at time zone 'Asia/Shanghai', 'YYYY-MM-DD') desc").
		Group("to_char(date at time zone 'Asia/Shanghai', 'YYYY-MM-DD')").
		Limit(limit).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	for i := range result {
		date, e := time.ParseInLocation(consts.DateLayout, result[i].Date, consts.DefaultLoc)
		if e != nil {
			return nil, e
		}
		result[i].Timestamp = consts.TransferTimeToTimestamp(date)
	}
	return result, nil
}

// 在Redis获取今天的访问量
func (d *DashboardDao) GetTodayPV(ctx context.Context) (response.DashboardTrends, error) {
	var result response.DashboardTrends
	today := consts.TransferTimeByLoc(time.Now()).Format(consts.DateLayout)
	// 调用函数获取今日PV
	uv, pv, _ := d.GetTodayPVUV(ctx)
	// 组装结果
	result.UV = uv
	result.PV = pv
	result.Date = today
	return result, nil
}

// 获取总日志数
func (d *DashboardDao) GetTotalLogs() (int64, error) {
	var result int64
	err := d.db.Table("visit_logs").
		Count(&result).
		Error
	if err != nil {
		return 0, err
	}
	return result, nil
}

// 获取访问的来源
func (d *DashboardDao) GetGeoDistribution(
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
	db := d.db.Table("visit_logs").
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
func (d *DashboardDao) GetErrorLogs(limit int) ([]response.ErrorLogItem, error) {
	var result []response.ErrorLogItem
	err := d.db.Table("visit_logs").
		Select("path,status,visit_time").
		Where("status != 200").
		Order("visit_time desc").
		Limit(limit).
		Scan(&result).
		Error
	if err != nil {
		return nil, err
	}
	for i := range result {
		result[i].Timestamp = consts.TransferTimeToTimestamp(result[i].Time)
	}
	return result, nil
}

// 获取前一天的PV和UV
func (d *DashboardDao) GetLastDayPVUV() (response.DashboardTrends, error) {
	yesterday := consts.TransferTimeByLoc(time.Now().AddDate(0, 0, -1)).Format(consts.DateLayout)
	// 复用结构体
	var result response.DashboardTrends
	err := d.db.Table("daily_article_stats").
		Select("coalesce(sum(pv), 0) as pv,coalesce(sum(uv), 0) as uv").
		Where("date = ?", yesterday).
		Scan(&result).
		Error

	if err != nil {
		return response.DashboardTrends{}, err
	}

	result.Date = yesterday
	return result, nil
}
