package utils

import "github.com/mssola/useragent"

func ParseUA(uaStr string) (device, os, browser string) {
	ua := useragent.New(uaStr)

	// 解析浏览器
	name, version := ua.Browser()
	browser = name + " " + version

	// 解析操作系统
	os = ua.OS()

	// 解析设备
	if ua.Mobile() {
		device = "mobile"
	} else if ua.Bot() {
		device = "bot"
	} else {
		device = "desktop"
	}
	return
}
