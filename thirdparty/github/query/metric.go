package query

import "github.com/shurcooL/githubv4"

type MetricQuery struct {
	Repository struct {
		Discussions struct {
			Nodes    []MetricDiscussion
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first,after: $after,orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type MetricReply struct {
	CreatedAt githubv4.DateTime
	Reactions struct {
		Nodes []MetricReaction
	} `graphql:"reactions(first:20)"`
}

type MetricReaction struct {
	CreatedAt githubv4.DateTime
}

type MetricComment struct {
	CreatedAt githubv4.DateTime
	Reactions struct {
		Nodes []MetricReaction
	} `graphql:"reactions(first:20)"`
	Replies struct {
		Nodes []MetricReply
	} `graphql:"replies(first:20)"`
}

type MetricDiscussion struct {
	UpdatedAt githubv4.DateTime
	Reactions struct {
		Nodes []MetricReaction
	} `graphql:"reactions(first:20)"`
	Comments struct {
		Nodes []MetricComment
	} `graphql:"comments(first:20)"`
}
