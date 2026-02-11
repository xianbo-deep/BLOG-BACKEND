package query

import (
	"github.com/shurcooL/githubv4"
)

/* 通用模块 */
// 查询结构体
type PageInfo struct {
	HasNextPage bool
	EndCursor   githubv4.String
}
