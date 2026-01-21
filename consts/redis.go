package consts

import (
	"fmt"
	"time"
)

const (
	// 分布式锁Key
	RedisLockKey = "cron:lock:sync"
	// Key前缀
	RedisKeyPrefix = "blog:stat:daily:"
	// 缓存前缀
	RedisCacheKeyPrefix = "blog:github:cache:"
	// 缓存版本号
	RedisGithubCacheVerKey = "blog:github:cache:ver"
)

const (
	RedisCacheKeySuffixMetric      = ":metric"
	RedisCacheKeySuffixTrend       = ":trend"
	RedisCacheKeySuffixActiveUsers = ":active_users"
	RedisCacheKeySuffixNewFeeds    = ":newfeeds"
)

const (

	// 在线人数
	RedisKeySuffixOnline = ":online"

	// 文章排名
	RedisKeySuffixPathRank = ":path_rank"

	// 总pv
	RedisKeySuffixTotalPV = ":total_pv"

	// 总uv
	RedisKeySuffixTotalUV = ":total_uv"

	// 页面访问次数
	RedisKeySuffixPathCount = ":latency:count"

	// 页面访问延迟
	RedisKeySuffixPathTotalLatency = ":latency:total"

	// 页面平均延迟
	RedisKeySuffixPathAvgLatency = ":latency:rank"
)

const (
	// 页面uv
	RedisKeyFmtPathUV = RedisKeyPrefix + "%s:uv:%s"
)

const (
	// Redis操作超时时间
	RedisOperationTimeout = 1 * TimeRangeSecond

	// Redis版本号过期时间
	RedisCacheVersionTimeout = 10 * TimeRangeMinute
)

// 获取Redis Key
func GetDailyStatKey(date string, suffix string) string {
	return RedisKeyPrefix + date + suffix
}

// 获取单页面 UV Key
func GetDailyPathUVKey(date string, path string) string {
	return fmt.Sprintf(RedisKeyFmtPathUV, date, path)
}

// 获取今日日期字符串
func GetTodayDate() string {
	loc, _ := time.LoadLocation(TimeLocation)
	return time.Now().In(loc).Format(DateLayout)
}

/* 获取github的缓存Key */
func GetGithubMetricCacheKey(ver int64, days int) string {
	return fmt.Sprintf("%s%d:%d%s", RedisCacheKeyPrefix, ver, days, RedisCacheKeySuffixMetric)
}

func GetGithubTrendCacheKey(ver int64, days int) string {
	return fmt.Sprintf("%s%d:%d%s", RedisCacheKeyPrefix, ver, days, RedisCacheKeySuffixTrend)
}

func GetGithubActiveUsersCacheKey(ver int64, limit int) string {
	return fmt.Sprintf("%s%d:%d%s", RedisCacheKeyPrefix, ver, limit, RedisCacheKeySuffixActiveUsers)
}

func GetGithubNewFeedsCacheKey(ver int64, limit int) string {
	return fmt.Sprintf("%s%d:%d%s", RedisCacheKeyPrefix, ver, limit, RedisCacheKeySuffixNewFeeds)
}
