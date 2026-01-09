package utils

import (
	"Blog-Backend/consts"

	refererparser "github.com/snowplow-referer-parser/golang-referer-parser"
)

func ParseReferer(raw string) (known bool, medium, source, term string) {
	if raw == "" {
		return false, consts.RefererDirect, "", ""
	}
	r := refererparser.Parse(raw)
	if !r.Known {
		// 无法知道来源
		return false, "unknown", "", ""
	}

	return true, r.Medium, r.Referer, r.SearchTerm
}
