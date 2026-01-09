package github

import (
	"Blog-Backend/consts"
	"context"
	"os"
	"sort"
	"time"

	"github.com/shurcooL/githubv4"
)

// TODO cnm 是人类做的吗 后续改成把结果存缓存里 注册个github的webhook监听 有动态就清缓存
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
func (s *DiscussionService) GetTrend(ctx context.Context, timeRangeDays int) ([]TrendItem, error) {
	// 预定义
	var after *githubv4.String
	cutoffTime := time.Now().AddDate(0, 0, -timeRangeDays)

	// 按日期统计
	trendMap := make(map[string]*TrendItem)

	for i := 0; i < timeRangeDays; i++ {
		date := time.Now().AddDate(0, 0, -i)
		key := date.Format("2006-01-02")
		trendMap[key] = &TrendItem{
			Date:           key,
			TotalReplies:   0,
			TotalComments:  0,
			TotalReactions: 0,
		}
	}

	for {
		var q struct {
			Repository struct {
				Discussion struct {
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
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
				} `graphql:"discussion(first: $first, after: $after)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		vars := map[string]interface{}{
			"first": consts.DefaultQuerySize,
			"after": after,
			"owner": s.owner,
			"repo":  s.repo,
		}

		err := s.github.Query(ctx, q, vars)
		if err != nil {
			return nil, err
		}

		for _, discussion := range q.Repository.Discussion.Nodes {
			for _, comment := range discussion.Comments.Nodes {
				if comment.CreatedAt.After(cutoffTime) {
					key := comment.CreatedAt.Time.Format("2006-01-02")
					trendMap[key].TotalComments++
				}
				for _, reply := range comment.Replies.Nodes {
					if reply.CreatedAt.After(cutoffTime) {
						key := reply.CreatedAt.Time.Format("2006-01-02")
						trendMap[key].TotalReplies++
					}
					for _, reaction := range reply.Reactions.Nodes {
						if reaction.CreatedAt.After(cutoffTime) {
							key := reaction.CreatedAt.Time.Format("2006-01-02")
							trendMap[key].TotalReactions++
						}
					}
				}
				for _, reaction := range comment.Reactions.Nodes {
					if reaction.CreatedAt.After(cutoffTime) {
						key := reaction.CreatedAt.Time.Format("2006-01-02")
						trendMap[key].TotalReactions++
					}
				}
			}
			for _, reaction := range discussion.Reactions.Nodes {
				if reaction.CreatedAt.After(cutoffTime) {
					key := reaction.CreatedAt.Time.Format("2006-01-02")
					trendMap[key].TotalReactions++
				}
			}
		}

		if !q.Repository.Discussion.PageInfo.HasNextPage {
			break
		}

		// 更新参数
		after = &q.Repository.Discussion.PageInfo.EndCursor
	}
	// 转换成切片
	trends := make([]TrendItem, 0, len(trendMap))
	for i := 0; i < timeRangeDays; i++ {
		date := time.Now().AddDate(0, 0, -i)
		key := date.Format("2006-01-02")
		trends = append(trends, *trendMap[key])
	}
	return trends, nil

}

// 返回活跃用户
func (s *DiscussionService) GetActiveUser(ctx context.Context, limit int) ([]ActiveUserItem, error) {
	var after *githubv4.String

	// 创建用户map
	userMap := make(map[string]*ActiveUserItem)
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
						} `graphql:"reactions(first:100)"`
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
										} `graphql:"reactions(first:100)"`
									}
								} `graphql:"replies(first:100)"`
								Reactions struct {
									Nodes []struct {
										User struct {
											AvatarUrl githubv4.String
											Url       githubv4.String
											Login     githubv4.String
										}
									}
								} `graphql:"reactions(first:100)"`
							}
						} `graphql:"comments(first:100)"`
					}
					PageInfo struct {
						HasNextPage bool
						EndCursor   githubv4.String
					}
				} `graphql:"discussion(first: $first, after: $after)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}
		vars := map[string]interface{}{
			"first": consts.DefaultQuerySize,
			"after": after,
			"owner": s.owner,
			"repo":  s.repo,
		}
		err := s.github.Query(ctx, q, vars)
		if err != nil {
			return nil, err
		}

		for _, discussion := range q.Repository.Discussions.Nodes {
			for _, comment := range discussion.Comments.Nodes {
				if comment.Author.Login != "" {
					login := string(comment.Author.Login)
					item, ok := userMap[login]
					if !ok {
						item = &ActiveUserItem{
							Name:       login,
							TotalFeeds: 0,
							Avatar:     string(comment.Author.AvatarUrl),
							URL:        string(comment.Author.Url),
						}
						userMap[login] = item
					}
					// 是指针，在这改会进行同步
					item.TotalFeeds++
				}
				for _, reaction := range comment.Reactions.Nodes {
					login := string(reaction.User.Login)
					if login == "" {
						continue
					}
					item, ok := userMap[login]
					if !ok {
						item = &ActiveUserItem{
							Name:       login,
							TotalFeeds: 0,
							Avatar:     string(reaction.User.AvatarUrl),
							URL:        string(reaction.User.Url),
						}
						userMap[login] = item
					}
					item.TotalFeeds++
				}
				for _, reply := range comment.Replies.Nodes {
					if reply.Author.Login != "" {
						login := string(reply.Author.Login)
						item, ok := userMap[login]
						if !ok {
							item = &ActiveUserItem{
								Name:       login,
								TotalFeeds: 0,
								Avatar:     string(reply.Author.AvatarUrl),
								URL:        string(reply.Author.Url),
							}
							userMap[login] = item
						}
						item.TotalFeeds++
					}
					for _, reaction := range reply.Reactions.Nodes {
						login := string(reaction.User.Login)
						if login == "" {
							continue
						}
						item, ok := userMap[login]
						if !ok {
							item = &ActiveUserItem{
								Name:       login,
								TotalFeeds: 0,
								Avatar:     string(reaction.User.AvatarUrl),
								URL:        string(reaction.User.Url),
							}
							userMap[login] = item
						}
						item.TotalFeeds++
					}
				}
			}

			for _, reaction := range discussion.Reactions.Nodes {
				login := string(reaction.User.Login)
				if login == "" {
					continue
				}
				item, ok := userMap[login]
				if !ok {
					item = &ActiveUserItem{
						Name:       login,
						TotalFeeds: 0,
						Avatar:     string(reaction.User.AvatarUrl),
						URL:        string(reaction.User.Url),
					}
					userMap[login] = item
				}
				item.TotalFeeds++
			}
		}

		if !q.Repository.Discussions.PageInfo.HasNextPage {
			break
		}

		after = &q.Repository.Discussions.PageInfo.EndCursor

	}

	// 进行数据组装返回
	activeusers := make([]ActiveUserItem, 0, len(userMap))
	for _, item := range userMap {
		activeusers = append(activeusers, *item)
	}

	// 按活跃度排序
	sort.Slice(activeusers, func(i, j int) bool {
		return activeusers[i].TotalFeeds > activeusers[j].TotalFeeds
	})

	if limit > 0 && limit < len(activeusers) {
		activeusers = activeusers[:limit]
	}

	return activeusers, nil
}
