package github

import (
	"Blog-Backend/consts"
	"context"
	"os"
	"sort"
	"time"

	"github.com/shurcooL/githubv4"
)

// TODO cnm 是人类做的吗 后续改成跑定时任务 反正也是给我自己看 如果频繁查询这也太慢了
type DiscussionService struct {
	github *githubv4.Client
	owner  string
	repo   string
}

func NewDiscussionService(github *githubv4.Client) *DiscussionService {
	return &DiscussionService{
		github: github,
		owner:  consts.OWNER,
		repo:   consts.REPO,
	}
}

// 返回三个指标
func (s *DiscussionService) GetTotalMetric(ctx context.Context, timeRangeDays int) (Metric, error) {
	// 声明为指针，返回nil
	var after *githubv4.String
	var cutoffTime time.Time
	totalComments := 0
	totalReplies := 0
	totalReactions := 0

	// 计算截止时间
	if timeRangeDays > 0 {
		cutoffTime = time.Now().AddDate(0, 0, -timeRangeDays)
	}

	// 准备discussion的排序参数
	orderBy := githubv4.DiscussionOrder{
		Field:     githubv4.DiscussionOrderFieldUpdatedAt,
		Direction: githubv4.OrderDirectionDesc,
	}
	for {
		var q struct {
			Repository struct {
				Discussions struct {
					Nodes []struct {
						UpdatedAt time.Time
						Reactions struct {
							Nodes []struct {
								CreatedAt time.Time
							}
						} `graphql:"reactions(first: 100)"`
						Comments struct {
							Nodes []struct {
								CreatedAt time.Time
								Reactions struct {
									Nodes []struct {
										CreatedAt time.Time
									}
								} `graphql:"reactions(first: 100)"`
								Replies struct {
									Nodes []struct {
										CreatedAt time.Time
										Reactions struct {
											Nodes []struct {
												CreatedAt time.Time
											} `graphql:"reactions(first: 100)"`
										}
									} `graphql:"nodes"`
								} `graphql:"replies(first: 100)"`
							}
						} `graphql:"comments(first: 100)"`
					}
					PageInfo struct {
						HasNextPage bool            // 是否还有下一页
						EndCursor   githubv4.String // 上一页的末尾光标
					}
				} `graphql:"discussions(first: $first,after: $after,orderBy: $orderBy)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		vars := map[string]interface{}{
			"owner":   githubv4.String(s.owner),
			"repo":    githubv4.String(s.repo),
			"after":   after,
			"first":   githubv4.Int(consts.DefaultQuerySize),
			"orderBy": orderBy,
		}

		err := s.github.Query(ctx, &q, vars)

		if err != nil {
			return Metric{}, err
		}
		shouldStop := false
		// 处理每一个discussion
		for _, discussion := range q.Repository.Discussions.Nodes {
			if timeRangeDays > 0 && discussion.UpdatedAt.Before(cutoffTime) {
				shouldStop = true
				break
			}
			// 累加discussion的评论
			for _, comment := range discussion.Comments.Nodes {
				if comment.CreatedAt.After(cutoffTime) {
					totalComments++
				}
				// 累加comment的回复
				for _, reply := range comment.Replies.Nodes {
					if reply.CreatedAt.After(cutoffTime) {
						totalReplies++
					}
					// 累加reply的回应
					for _, reaction := range reply.Reactions.Nodes {
						if reaction.CreatedAt.After(cutoffTime) {
							totalReactions++
						}
					}
				}
				// 累加comment的回应
				for _, reaction := range comment.Reactions.Nodes {
					if reaction.CreatedAt.After(cutoffTime) {
						totalReactions++
					}
				}
			}
			// 累加discussion的回应
			for _, reaction := range discussion.Reactions.Nodes {
				if reaction.CreatedAt.After(cutoffTime) {
					totalReactions++
				}
			}
		}

		// 若discussion更新时间早于截止时间或者没有下一页
		if shouldStop || !q.Repository.Discussions.PageInfo.HasNextPage {
			break
		}

		// 更新游标
		after = &q.Repository.Discussions.PageInfo.EndCursor
	}

	return Metric{
		TotalComments:  int64(totalComments),
		TotalReplies:   int64(totalReplies),
		TotalReactions: int64(totalReactions),
	}, nil

}

// 返回最新动态
func (s *DiscussionService) GetNewFeed(ctx context.Context, limit int) ([]NewFeedItem, error) {
	var after *githubv4.String
	var allItems []NewFeedItem
	// 不统计7天前的
	cutoffTime := time.Now().AddDate(0, 0, -consts.TimeRangeWeek)
	for {
		var q struct {
			Repository struct {
				Discussions struct {
					Nodes []struct {
						Title     githubv4.String
						UpdatedAt time.Time
						Reactions struct {
							Nodes []struct {
								CreatedAt time.Time
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
								CreatedAt time.Time
								Reactions struct {
									Nodes []struct {
										CreatedAt time.Time
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
										CreatedAt time.Time
										BodyText  githubv4.String
										Reactions struct {
											Nodes []struct {
												CreatedAt time.Time
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
					PageInfo struct {
						HasNextPage bool            // 是否还有下一页
						EndCursor   githubv4.String // 上一页的末尾光标
					}
				} `graphql:"discussions(first: $first,after: $after,orderBy:$orderBy)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}
		vars := map[string]interface{}{
			"owner": githubv4.String(s.owner),
			"repo":  githubv4.String(s.repo),
			"orderBy": githubv4.DiscussionOrder{
				Field:     githubv4.DiscussionOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			"first": githubv4.Int(consts.DefaultQuerySize),
			"after": after,
		}

		err := s.github.Query(ctx, &q, vars)

		if err != nil {
			return nil, err
		}

		// 是否应该暂停
		shouldStop := false

		for _, discussion := range q.Repository.Discussions.Nodes {
			if discussion.UpdatedAt.Before(cutoffTime) {
				shouldStop = true
				break
			}
			// 处理reaction
			for _, reaction := range discussion.Reactions.Nodes {
				if reaction.CreatedAt.After(cutoffTime) {
					allItems = append(allItems, NewFeedItem{
						EventType: consts.Reaction,
						Name:      string(reaction.User.Login),
						Path:      os.Getenv(consts.EnvBaseURL) + string(discussion.Title),
						Content:   string(reaction.Content),
						Avatar:    string(reaction.User.AvatarUrl),
						Time:      reaction.CreatedAt.Format(time.RFC3339),
						URL:       string(reaction.User.Url),
					})
				}
			}
			// 处理评论
			for _, comment := range discussion.Comments.Nodes {
				if comment.CreatedAt.After(cutoffTime) {
					allItems = append(allItems, NewFeedItem{
						EventType: consts.Comment,
						Name:      string(comment.Author.Login),
						Path:      os.Getenv(consts.EnvBaseURL) + string(discussion.Title),
						Content:   string(comment.BodyText),
						Avatar:    string(comment.Author.AvatarUrl),
						Time:      comment.CreatedAt.Format(time.RFC3339),
						URL:       string(comment.Author.Url),
					})
				}
				// 处理评论的reaction
				for _, reaction := range comment.Reactions.Nodes {
					if reaction.CreatedAt.After(cutoffTime) {
						allItems = append(allItems, NewFeedItem{
							EventType: consts.Reaction,
							Name:      string(reaction.User.Login),
							Path:      os.Getenv(consts.EnvBaseURL) + string(discussion.Title),
							Content:   string(reaction.Content),
							Avatar:    string(reaction.User.AvatarUrl),
							Time:      reaction.CreatedAt.Format(time.RFC3339),
							URL:       string(reaction.User.Url),
						})
					}
				}

				// 处理评论的回复
				for _, reply := range comment.Replies.Nodes {
					if reply.CreatedAt.After(cutoffTime) {
						allItems = append(allItems, NewFeedItem{
							EventType:      consts.Reply,
							Name:           string(reply.Author.Login),
							Path:           os.Getenv(consts.EnvBaseURL) + string(discussion.Title),
							Content:        string(reply.BodyText),
							Avatar:         string(reply.Author.AvatarUrl),
							Time:           reply.CreatedAt.Format(time.RFC3339),
							URL:            string(reply.Author.Url),
							ReplyToName:    string(comment.Author.Login),
							ReplyToAvatar:  string(comment.Author.AvatarUrl),
							ReplyToContent: string(comment.BodyText),
						})
					}
					// 处理回复的reaction
					for _, reaction := range reply.Reactions.Nodes {
						if reaction.CreatedAt.After(cutoffTime) {
							allItems = append(allItems, NewFeedItem{
								EventType: consts.Reaction,
								Name:      string(reaction.User.Login),
								Path:      os.Getenv(consts.EnvBaseURL) + string(discussion.Title),
								Content:   string(reaction.Content),
								Avatar:    string(reaction.User.AvatarUrl),
								Time:      reaction.CreatedAt.Format(time.RFC3339),
								URL:       string(reaction.User.Url),
							})
						}
					}
				}
			}
		}
		// 若长度过长
		if len(allItems) > 2*consts.DefaultQuerySize {
			break
		}
		// 退出 装载信息
		if shouldStop || !q.Repository.Discussions.PageInfo.HasNextPage {
			break
		}

		// 更新after
		after = &q.Repository.Discussions.PageInfo.EndCursor
	}

	// 对动态按照时间进行排序
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].Time > allItems[j].Time
	})

	// 截取动态
	if len(allItems) > limit {
		allItems = allItems[:limit]
	}

	return allItems, nil
}

// 返回趋势
func (s *DiscussionService) GetTrend(ctx context.Context) ([]TrendItem, error) {
	var after *githubv4.String
	var trends []TrendItem
	cutoffTime := time.Now().AddDate(0, 0, -consts.TimeRangeWeek)
	for {
		var q struct {
			Repository struct {
				Discussion struct {
					Nodes []struct {
						Reactions struct {
							Nodes []struct {
								CreatedAt githubv4.DateTime
							}
						}
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
										}
									}
								}
								Reactions struct {
									Nodes []struct {
										CreatedAt githubv4.DateTime
									}
								}
							}
						}
					}
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
				} `graphql:"discussion(first: $first, after: $after)"`
			} `graphql:"repository(owner: $owner, name: $name)"`
		}
	}
}

// 返回活跃用户
func (s *DiscussionService) GetActiveUser(ctx context.Context) ([]ActiveUserItem, error) {
	var after *githubv4.String
	var activeUsers []ActiveUserItem
	for {
		var q struct {
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
						}
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
										}
									}
								}
								Reactions struct {
									Nodes []struct {
										User struct {
											AvatarUrl githubv4.String
											Url       githubv4.String
											Login     githubv4.String
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
