package github

import (
	"time"

	"github.com/shurcooL/graphql"
)

// 查询时间节点
type CreatedAtNode struct {
	CreatedAt time.Time
}

// Reaction集合
type ReactionConnection struct {
	Nodes []CreatedAtNode `graphql:"nodes"`
}

// Reply节点
type ReplyNode struct {
	CreatedAt time.Time
	Reactions ReactionConnection `graphql:"reactions(first: 100)"`
}

// Comment节点
type CommentNode struct {
	CreatedAt time.Time
	Reactions []ReactionConnection `graphql:"reactions(first: 100)"`
	Replies   struct {
		Nodes []ReplyNode `graphql:"nodes"`
	} `graphql:"replies(last: 100)"`
}

// Discussion节点
type DiscussionNode struct {
	UpdatedAt time.Time
	Reactions []ReactionConnection `graphql:"reactions(first: 100)"`
	Comments  struct {
		Nodes []CommentNode `graphql:"nodes"`
	} `graphql:"comments(first: 100)"`
}

// 主体结构
type DiscussionQuery struct {
	Repository struct {
		Disscussion struct {
			Nodes    []DiscussionNode `graphql:"nodes"`
			PageInfo struct {
				HasNextPage bool
				EndCursor   graphql.String
			}
		} `graphql:"discussions(first: 100,after: $after,orderBy:$orderBy)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}
