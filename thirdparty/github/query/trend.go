package query

import "github.com/shurcooL/githubv4"

type TrendQuery struct {
	Repository struct {
		Discussions struct {
			Nodes    []TrendDiscussion
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first,after: $after)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type TrendReply struct {
	CreatedAt githubv4.DateTime
	Reactions struct {
		Nodes []TrendReaction
	} `graphql:"reactions(first:20)"`
}

type TrendReaction struct {
	CreatedAt githubv4.DateTime
}

type TrendComment struct {
	CreatedAt githubv4.DateTime
	Replies   struct {
		Nodes []TrendReply
	} `graphql:"replies(first:20)"`
	Reactions struct {
		Nodes []TrendReaction
	} `graphql:"reactions(first:20)"`
}

type TrendDiscussion struct {
	Reactions struct {
		Nodes []TrendReaction
	} `graphql:"reactions(first:20)"`
	Comments struct {
		Nodes []TrendComment
	} `graphql:"comments(first:20)"`
}
