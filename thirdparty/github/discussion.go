package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"context"
	"go/ast"
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
				} `graphql:"discussions(first: 100,after: $after,orderBy: $orderBy)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		vars := map[string]interface{}{
			"owner":   githubv4.String(s.owner),
			"repo":    githubv4.String(s.repo),
			"after":   after,
			"orderBy": orderBy,
		}

		err := s.github.Query(ctx, &q, vars)

		if err != nil {
			return response.Metric{}, err
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

	return response.Metric{
		TotalComments:  int64(totalComments),
		TotalReplies:   int64(totalReplies),
		TotalReactions: int64(totalReactions),
	}, nil

}

// 返回最新动态
func (s *DiscussionService) GetNewFeed(ctx context.Context) ([]response.NewFeedItem, error) {

}

// 返回趋势
func (s *DiscussionService) GetTrend(ctx context.Context) ([]response.TrendItem, error) {

}

// 返回活跃用户
func (s *DiscussionService) GetActiveUser(ctx context.Context) ([]response.ActiveUserItem, error) {

}
