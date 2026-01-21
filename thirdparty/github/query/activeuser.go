package query

import "github.com/shurcooL/githubv4"

type ActiveUserQuery struct {
	Repository struct {
		Discussions struct {
			Nodes []struct {
				Reactions struct {
					Nodes []struct {
						User struct {
							AvatarUrl githubv4.String
							Url       githubv4.String
							Login     githubv4.String
						}
					}
				} `graphql:"reactions(first:20)"`
				Comments struct {
					Nodes []struct {
						Author struct {
							AvatarUrl githubv4.String
							Url       githubv4.String
							Login     githubv4.String
						}
						Replies struct {
							Nodes []struct {
								Author struct {
									AvatarUrl githubv4.String
									Url       githubv4.String
									Login     githubv4.String
								}
								Reactions struct {
									Nodes []struct {
										User struct {
											AvatarUrl githubv4.String
											Url       githubv4.String
											Login     githubv4.String
										}
									}
								} `graphql:"reactions(first:20)"`
							}
						} `graphql:"replies(first:20)"`
						Reactions struct {
							Nodes []struct {
								User struct {
									AvatarUrl githubv4.String
									Url       githubv4.String
									Login     githubv4.String
								}
							}
						} `graphql:"reactions(first:20)"`
					}
				} `graphql:"comments(first:20)"`
			}
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"discussions(first: $first, after: $after)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}
