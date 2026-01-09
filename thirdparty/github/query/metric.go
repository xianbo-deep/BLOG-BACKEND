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
				} `graphql:"reactions(first: 100)"`
				Comments struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
						Reactions struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
							}
						} `graphql:"reactions(first: 100)"`
						Replies struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								Reactions struct {
									Nodes []struct {
										CreatedAt githubv4.DateTime
									} `graphql:"reactions(first: 100)"`
								}
							} `graphql:"nodes"`
						} `graphql:"replies(first: 100)"`
					}
				} `graphql:"comments(first: 100)"`
			}
			PageInfo
		} `graphql:"discussions(first: $first,after: $after,orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}
