package service

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/thirdparty/github/query"
	"context"
	"os"
	"time"

	"github.com/shurcooL/githubv4"
)

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
func (s *DiscussionService) GetTotalMetric(ctx context.Context, timeRangeDays int) (response.Metric, error) {
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
		var q query.MetricQuery

		vars := map[string]interface{}{
			"owner":   githubv4.String(s.owner),
			"repo":    githubv4.String(s.repo),
			"after":   after,
			"first":   githubv4.Int(consts.DefaultQuerySize),
			"orderBy": orderBy,
		}

		err := s.github.Query(ctx, &q, vars)

		if err != nil {
			return response.Metric{}, err
		}
		shouldStop := false
		// 处理每一个discussion
		for _, discussion := range q.Repository.Discussions.Nodes {
			if timeRangeDays > 0 && !shouldCount(cutoffTime, discussion.UpdatedAt.Time) {
				shouldStop = true
				break
			}
			// 累加discussion的评论
			for _, comment := range discussion.Comments.Nodes {
				if shouldCount(cutoffTime, comment.CreatedAt.Time) {
					totalComments++
				}
				// 累加comment的回复
				for _, reply := range comment.Replies.Nodes {
					if shouldCount(cutoffTime, reply.CreatedAt.Time) {
						totalReplies++
					}
					// 累加reply的回应
					for _, reaction := range reply.Reactions.Nodes {
						if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
							totalReactions++
						}
					}
				}
				// 累加comment的回应
				for _, reaction := range comment.Reactions.Nodes {
					if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
						totalReactions++
					}
				}
			}
			// 累加discussion的回应
			for _, reaction := range discussion.Reactions.Nodes {
				if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
					totalReactions++
				}
			}
		}

		// 若discussion更新时间早于截止时间或者没有下一页
		if shouldStop || !q.Repository.Discussions.PageInfo.HasNextPage {
			break
		}

		// 更新游标
		after = nextCursor(q.Repository.Discussions.PageInfo)
	}

	return response.Metric{
		TotalComments:  int64(totalComments),
		TotalReplies:   int64(totalReplies),
		TotalReactions: int64(totalReactions),
	}, nil

}

// 返回最新动态
func (s *DiscussionService) GetNewFeed(ctx context.Context, limit int) ([]*response.NewFeedItem, error) {
	var after *githubv4.String
	var allItems []*response.NewFeedItem
	// 不统计7天前的
	cutoffTime := time.Now().AddDate(0, 0, -consts.Week)
	for {
		var q query.FeedQuery

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
			if !shouldCount(cutoffTime, discussion.UpdatedAt.Time) {
				shouldStop = true
				break
			}
			// 处理reaction
			for _, reaction := range discussion.Reactions.Nodes {
				if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
					allItems = append(allItems, &response.NewFeedItem{
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
				if shouldCount(cutoffTime, comment.CreatedAt.Time) {
					allItems = append(allItems, &response.NewFeedItem{
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
					if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
						allItems = append(allItems, &response.NewFeedItem{
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
					if shouldCount(cutoffTime, reply.CreatedAt.Time) {
						allItems = append(allItems, &response.NewFeedItem{
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
						if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
							allItems = append(allItems, &response.NewFeedItem{
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
		after = nextCursor(q.Repository.Discussions.PageInfo)
	}

	return handleNewFeedRes(allItems, limit)
}

// 返回趋势
func (s *DiscussionService) GetTrend(ctx context.Context, timeRangeDays int) ([]response.TrendItem, error) {
	// 预定义
	var after *githubv4.String
	cutoffTime := time.Now().AddDate(0, 0, -timeRangeDays)

	// 按日期统计
	trendMap := make(map[string]*response.TrendItem)

	for i := 0; i < timeRangeDays; i++ {
		date := time.Now().AddDate(0, 0, -i)
		key := date.Format(consts.DateLayout)
		trendMap[key] = &response.TrendItem{
			Date:           key,
			TotalReplies:   0,
			TotalComments:  0,
			TotalReactions: 0,
		}
	}

	for {
		var q query.TrendQuery

		vars := map[string]interface{}{
			"first": githubv4.Int(consts.DefaultQuerySize),
			"after": after,
			"owner": s.owner,
			"repo":  s.repo,
		}

		err := s.github.Query(ctx, &q, vars)
		if err != nil {
			return nil, err
		}

		for _, discussion := range q.Repository.Discussions.Nodes {
			for _, comment := range discussion.Comments.Nodes {
				if shouldCount(cutoffTime, comment.CreatedAt.Time) {
					addTrendItemComment(trendMap, comment.CreatedAt.Time)
				}
				for _, reply := range comment.Replies.Nodes {
					if shouldCount(cutoffTime, reply.CreatedAt.Time) {
						addTrendItemReply(trendMap, reply.CreatedAt.Time)
					}
					for _, reaction := range reply.Reactions.Nodes {
						if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
							addTrendItemReaction(trendMap, reaction.CreatedAt.Time)
						}
					}
				}
				for _, reaction := range comment.Reactions.Nodes {
					if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
						addTrendItemReaction(trendMap, reaction.CreatedAt.Time)
					}
				}
			}
			for _, reaction := range discussion.Reactions.Nodes {
				if shouldCount(cutoffTime, reaction.CreatedAt.Time) {
					addTrendItemReaction(trendMap, reaction.CreatedAt.Time)
				}
			}
		}

		if !q.Repository.Discussions.PageInfo.HasNextPage {
			break
		}

		// 更新参数
		after = nextCursor(q.Repository.Discussions.PageInfo)
	}

	return handleTrendRes(trendMap, timeRangeDays)
}

// 返回活跃用户
func (s *DiscussionService) GetActiveUser(ctx context.Context, limit int) ([]response.ActiveUserItem, error) {
	var after *githubv4.String

	// 创建用户map
	userMap := make(map[string]*response.ActiveUserItem)
	for {
		var q query.ActiveUserQuery

		vars := map[string]interface{}{
			"first": githubv4.Int(consts.DefaultQuerySize),
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
				addActiveUser(userMap, string(comment.Author.Login), string(comment.Author.AvatarUrl), string(comment.Author.Url))
				for _, reaction := range comment.Reactions.Nodes {
					addActiveUser(userMap, string(reaction.User.Login), string(reaction.User.AvatarUrl), string(reaction.User.Url))
				}
				for _, reply := range comment.Replies.Nodes {
					addActiveUser(userMap, string(reply.Author.Login), string(reply.Author.AvatarUrl), string(reply.Author.Url))
					for _, reaction := range reply.Reactions.Nodes {
						addActiveUser(userMap, string(reaction.User.Login), string(reaction.User.AvatarUrl), string(reaction.User.Url))
					}
				}
			}

			for _, reaction := range discussion.Reactions.Nodes {
				addActiveUser(userMap, string(reaction.User.Login), string(reaction.User.AvatarUrl), string(reaction.User.Url))
			}
		}

		if !q.Repository.Discussions.PageInfo.HasNextPage {
			break
		}

		// 更新页数
		after = nextCursor(q.Repository.Discussions.PageInfo)

	}

	return handleActiveUserRes(userMap, limit)
}
