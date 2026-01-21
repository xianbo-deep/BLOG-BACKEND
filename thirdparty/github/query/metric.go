package query

import "github.com/shurcooL/githubv4"

type MetricQuery struct {
	Repository struct {
		Discussions struct {
			Nodes []struct {
				UpdatedAt githubv4.DateTime
				Reactions struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
					}
				} `graphql:"reactions(first:20)"`
				Comments struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
						Reactions struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
							}
						} `graphql:"reactions(first:20)"`
						Replies struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								Reactions struct {
									Nodes []struct {
										CreatedAt githubv4.DateTime
									}
								} `graphql:"reactions(first:20)"`
							} `graphql:"nodes"`
						} `graphql:"replies(first:20)"`
					}
				} `graphql:"comments(first:20)"`
			}
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first,after: $after,orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}
