package utils

import (
	"github.com/ua-parser/uap-go/uaparser"
)

func ParseUA(uaStr string) (device, os, browser string) {
	parser, _ := uaparser.New()

	client := parser.Parse(uaStr)
	// 解析浏览器
	browser = client.UserAgent.Family

	// 解析操作系统
	os = client.Os.Family

	// 解析设备
	if client.Device.Family == "iPhone" || client.Device.Family == "Android" {
		device = "mobile"
	} else if client.Device.Family == "Bot" {
		device = "bot"
	} else {
		device = "desktop"
	}
	return
}
