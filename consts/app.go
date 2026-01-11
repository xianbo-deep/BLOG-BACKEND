package consts

import "time"

const (
	// 日期格式
	// TODO 记得看看为什么
	DateLayout = "2006-01-02"

	// 标准时间格式
	TimeLayout = "2006-01-02 15:04:05"

	// 默认起始页数
	DefaultPage = 1

	// 默认分页大小
	DefaultPageSize = 20

	// 最大分页大小
	MaxPageSize = 100

	// JWT过期时间
	JwtTokenExpireDuration = time.Hour * 24

	// JWT 签发者
	JwtIssuer = "xbZhong"

	// 缓存过期时间
	CacheExpireDuration = 24 * time.Hour

	// 常用的时间
	TimeRangeHour = time.Hour * 1

	TimeRangeDay = time.Hour * 24

	TimeRangeWeek = time.Hour * 24 * 7

	TimeRangeMonth = time.Hour * 24 * 30

	TimeRangeYear = time.Hour * 24 * 365
)
