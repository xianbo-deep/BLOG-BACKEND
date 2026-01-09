package utils

import (
	"Blog-Backend/consts"

	refererparser "github.com/snowplow-referer-parser/golang-referer-parser"
)

func ParseReferer(raw string) (known bool, medium, source string) {
	if raw == "" {
		return false, consts.RefererDirect, ""
	}
	r := refererparser.Parse(raw)
	if !r.Known {
		// 无法知道来源
		return false, consts.RefererUnknown, ""
	}

	// 返回媒介类型、具体来源、搜索关键词
	return true, r.Medium, r.Referer
}
