package query

import "github.com/shurcooL/githubv4"

type DiscussionDigestQuery struct {
	Repository struct {
		Discussions struct {
			Nodes []struct {
				Title     githubv4.String
				UpdatedAt githubv4.DateTime
				Reactions struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
						AvatarUrl githubv4.String
						Url       githubv4.String
						Login     githubv4.String
						Content   githubv4.ReactionContent
					}
				} `graphql:"reactions(first: 20)"`
				Comments struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
						BodyText  githubv4.String
						AvatarUrl githubv4.String
						Url       githubv4.String
						Login     githubv4.String
						Reactions struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								AvatarUrl githubv4.String
								Url       githubv4.String
								Login     githubv4.String
								Content   githubv4.ReactionContent
							}
						} `graphql:"reactions(first: 20)"`
						Replies struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								AvatarUrl githubv4.String
								Url       githubv4.String
								Login     githubv4.String
								BodyText  githubv4.String
								Reactions struct {
									Nodes []struct {
										CreatedAt githubv4.DateTime
										AvatarUrl githubv4.String
										Url       githubv4.String
										Login     githubv4.String
										Content   githubv4.ReactionContent
									}
								} `graphql:"reactions(first: 20)"`
							}
						} `graphql:"replies(first: 20)"`
					}
				} `graphql:"comments(first: 20)"`
			}
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first, after: $after)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type DiscussionDigestAuthor struct{}
