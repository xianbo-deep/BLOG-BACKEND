package query

import "github.com/shurcooL/githubv4"

type TrendQuery struct {
	Repository struct {
		Discussions struct {
			Nodes []struct {
				Reactions struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
					}
				} `graphql:"reactions(first:100)"`
				Comments struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
						Replies   struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								Reactions struct {
									Nodes []struct {
										CreatedAt githubv4.DateTime
									}
								} `graphql:"reactions(first:100)"`
							}
						} `graphql:"replies(first:100)"`
						Reactions struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
							}
						} `graphql:"reactions(first:100)"`
					}
				} `graphql:"comments(first:100)"`
			}
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first,after: $after)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}
