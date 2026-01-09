package query

import "github.com/shurcooL/githubv4"

type FeedQuery struct {
	Repository struct {
		Discussions struct {
			Nodes []struct {
				Title     githubv4.String
				UpdatedAt githubv4.DateTime
				Reactions struct {
					Nodes []struct {
						CreatedAt githubv4.DateTime
						User      struct {
							AvatarUrl githubv4.String
							Url       githubv4.String
							Login     githubv4.String
						}
						Content githubv4.ReactionContent
					}
				} `graphql:"reactions(last: 20)"`
				Comments struct {
					Nodes []struct {
						BodyText  githubv4.String
						CreatedAt githubv4.DateTime
						Reactions struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								User      struct {
									AvatarUrl githubv4.String
									Url       githubv4.String
									Login     githubv4.String
								}
								Content githubv4.ReactionContent
							}
						} `graphql:"reactions(last: 20)"`
						Replies struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
								BodyText  githubv4.String
								Reactions struct {
									Nodes []struct {
										CreatedAt githubv4.DateTime
										User      struct {
											AvatarUrl githubv4.String
											Url       githubv4.String
											Login     githubv4.String
										}
										Content githubv4.ReactionContent
									}
								} `graphql:"reactions(last: 20)"`
								Author struct {
									AvatarUrl githubv4.String
									Url       githubv4.String
									Login     githubv4.String
								}
							}
						} `graphql:"replies(last: 20)"`
						Author struct {
							AvatarUrl githubv4.String
							Url       githubv4.String
							Login     githubv4.String
						}
					}
				} `graphql:"comments(last: 20)"`
			}
			PageInfo
		} `graphql:"discussions(first: $first,after: $after,orderBy:$orderBy)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}
