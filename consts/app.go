package consts

import (
	"context"
	"time"
)

const (
	RequestMetaKey = "requestMeta"
)

const (
	// 日期格式
	DateLayout = "2006-01-02"

	// 标准时间格式
	TimeLayout = "2006-01-02 15:04:05"

	// 时区
	TimeLocation = "Asia/Shanghai"
)

const (
	// 默认起始页数
	DefaultPage = 1

	// 默认分页大小
	DefaultPageSize = 20

	// 最大分页大小
	MaxPageSize = 100
)

const (
	// JWT过期时间
	JwtTokenExpireDuration = time.Hour * 24

	// JWT 签发者
	JwtIssuer = "xbZhong"

	// 缓存过期时间
	CacheExpireDuration = 24 * time.Hour
)

const (
	// 常用的时间
	TimeRangeSecond = time.Second

	TimeRangeMinute = time.Minute

	TimeRangeHour = time.Hour * 1

	TimeRangeDay = time.Hour * 24

	TimeRangeWeek = time.Hour * 24 * 7

	TimeRangeMonth = time.Hour * 24 * 30

	TimeRangeYear = time.Hour * 24 * 365
)

const (
	// 请求超时时间
	RequestTimeout = 2 * TimeRangeSecond
)

var (
	// 声明时区
	DefaultLoc = mustLoadLocation(TimeLocation)
)

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}

func GetCurrentUTCTime() time.Time {
	return time.Now().UTC()
}

func TransferTimeByLoc(t time.Time) time.Time {
	return t.In(DefaultLoc)
}

func TransferTimeToTimestamp(t time.Time) int64 {
	return t.In(DefaultLoc).UnixMilli()
}

func GetTimeoutContext(ctx context.Context, time time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time)
}
