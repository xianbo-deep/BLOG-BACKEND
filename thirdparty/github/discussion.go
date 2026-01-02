package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"context"
	"go/ast"

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
func (s *DiscussionService) GetTotalMertic(ctx context.Context) (response.Metric, error) {
	var after githubv4.String
	totalComments := 0
	totalReplies := 0
	totalResponses := 0
	for {
		var q struct {
			Repository struct {
				Discussions struct {
					TotalCount int
					Nodes      []struct {
						Reactions struct {
							TotalCount int
						}
						Comments struct {
							TotalCount int
							Nodes      []struct {
								Reactions struct {
									TotalCount int
								}
								Replies struct {
									TotalCount int
									Nodes      []struct {
										Reactions struct {
											TotalCount int
										}
									} `graphql:"nodes"`
								}
							}
						} `graphql:"comments(first: 50)"`
					}
					PageInfo struct {
						HasNextPage bool            // 是否还有下一页
						EndCursor   githubv4.String // 上一页的末尾光标
					}
				} `graphql:"discussions(first: 50,after: $after)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		vars := map[string]interface{}{
			"owner": githubv4.String(s.owner),
			"repo":  githubv4.String(s.repo),
			"after": after,
		}

		err := s.github.Query(ctx, &q, vars)

		if err != nil {
			return nil, err
		}

	}

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
