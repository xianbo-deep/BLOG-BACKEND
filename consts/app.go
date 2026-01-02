package consts

import "time"

const (
	// 日期格式
	DateLayout = "2006-01-02"

	// 标准时间格式
	TimeLayout = "2006-01-02 15:04:05"

	// 默认分页大小
	DefaultPageSize = 10

	// 最大分页大小
	MaxPageSize = 100

	// JWT过期时间
	JwtTokenExpireDuration = time.Hour * 24

	// JWT 签发者
	JwtIssuer = "xbZhong"
)
