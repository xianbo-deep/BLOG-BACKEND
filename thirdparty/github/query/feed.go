package query

import "github.com/shurcooL/githubv4"

type FeedQuery struct {
	Repository struct {
		Discussions struct {
			Nodes    []FeedDiscussion
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first,after: $after,orderBy:$orderBy)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

/* 最新动态 */
type FeedUser struct {
	AvatarUrl githubv4.String
	Url       githubv4.String
	Login     githubv4.String
}

type FeedReaction struct {
	CreatedAt githubv4.DateTime
	User      FeedUser
	Content   githubv4.ReactionContent
}

type FeedReply struct {
	CreatedAt githubv4.DateTime
	BodyText  githubv4.String
	Reactions struct {
		Nodes []FeedReaction
	} `graphql:"reactions(first: 20)"`
	Author FeedUser
}

type FeedComment struct {
	BodyText  githubv4.String
	CreatedAt githubv4.DateTime
	Reactions struct {
		Nodes []FeedReaction
	} `graphql:"reactions(first: 20)"`
	Replies struct {
		Nodes []FeedReply
	} `graphql:"replies(first: 20)"`
	Author FeedUser
}

type FeedDiscussion struct {
	Title     githubv4.String
	UpdatedAt githubv4.DateTime
	Reactions struct {
		Nodes []FeedReaction
	} `graphql:"reactions(first: 20)"`
	Comments struct {
		Nodes []FeedComment
	} `graphql:"comments(first: 20)"`
}
