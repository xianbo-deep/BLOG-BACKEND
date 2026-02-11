package query

import "github.com/shurcooL/githubv4"

type ActiveUserQuery struct {
	Repository struct {
		Discussions struct {
			Nodes    []ActiveUserDiscussion
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first, after: $after)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type ActiveUserAuthor struct {
	AvatarUrl githubv4.String
	Url       githubv4.String
	Login     githubv4.String
}

type ActiveUserReaction struct {
	User ActiveUserAuthor
}

type ActiveUserReply struct {
	Author    ActiveUserAuthor
	Reactions struct {
		Nodes []ActiveUserReaction
	} `graphql:"reactions(first:20)"`
}

type ActiveUserComment struct {
	Author  ActiveUserAuthor
	Replies struct {
		Nodes []ActiveUserReply
	} `graphql:"replies(first:20)"`
	Reactions struct {
		Nodes []ActiveUserReaction
	} `graphql:"reactions(first:20)"`
}

type ActiveUserDiscussion struct {
	Reactions struct {
		Nodes []ActiveUserReaction
	} `graphql:"reactions(first:20)"`
	Comments struct {
		Nodes []ActiveUserComment
	} `graphql:"comments(first:20)"`
}
