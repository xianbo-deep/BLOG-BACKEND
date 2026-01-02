package consts

import (
	"fmt"
	"time"
)

const (
	RedisKeyPrefix = "blog:stat:daily:"

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

	// 页面uv
	RedisKeyFmtPathUV = RedisKeyPrefix + "%s:uv:%s"
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
	return time.Now().Format(DateLayout)
}
