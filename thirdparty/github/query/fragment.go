package query

import (
	"github.com/shurcooL/githubv4"
)

// TODO 后续要把大的查询结构体的定义进行优化
// 查询结构体
type PageInfo struct {
	HasNextPage bool
	EndCursor   githubv4.String
}
