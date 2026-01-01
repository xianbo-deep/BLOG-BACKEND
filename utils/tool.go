package utils

import "time"

// 获取今日日期字符串，格式：20060102
func getTodayDate() string {
	return time.Now().Format("20060102")
}
