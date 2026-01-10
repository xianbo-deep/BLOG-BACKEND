package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao/cache"
	"Blog-Backend/thirdparty/github"
	"Blog-Backend/thirdparty/github/service"
	"context"
)

type CommentService struct {
	cached *cache.CacheDAO
	ds     *service.DiscussionService
}

func NewCommentService() *CommentService {
	return &CommentService{cached: cache.NewCacheDAO(), ds: service.NewDiscussionService(github.NewClient())}
}

func (s *CommentService) GetDiscussionMetric(c context.Context, days int) (response.Metric, error) {
	key := consts.GetGithubMetricCacheKey(days)
	var cached response.Metric
	ok, err := s.cached.GetJSON(c, key, &cached)
	if err != nil {
		return response.Metric{}, err
	}
	if ok {
		return cached, nil
	}

	res, err := s.ds.GetTotalMetric(c, days)
	if err != nil {
		return response.Metric{}, err
	}
	_ = s.cached.SetJSON(c, key, res, consts.CacheExpireDuration)
	return res, nil
}

func (s *CommentService) GetDiscussionTrend(c context.Context, days int) ([]response.TrendItem, error) {
	key := consts.GetGithubTrendCacheKey(days)
	var cached []response.TrendItem
	ok, err := s.cached.GetJSON(c, key, &cached)
	if err != nil {
		return []response.TrendItem{}, err
	}
	if ok {
		return cached, nil
	}
	res, err := s.ds.GetTrend(c, days)
	if err != nil {
		return []response.TrendItem{}, err
	}
	_ = s.cached.SetJSON(c, key, res, consts.CacheExpireDuration)
	return res, nil
}

func (s *CommentService) GetDiscussionNewFeed(c context.Context, limit int) ([]*response.NewFeedItem, error) {
	key := consts.GetGithubNewFeedsCacheKey(limit)
	var cached []*response.NewFeedItem
	ok, err := s.cached.GetJSON(c, key, &cached)
	if err != nil {
		return nil, err
	}
	if ok {
		return cached, nil
	}
	res, err := s.ds.GetNewFeed(c, limit)
	if err != nil {
		return nil, err
	}
	_ = s.cached.SetJSON(c, key, res, consts.CacheExpireDuration)
	return res, nil
}

func (s *CommentService) GetDiscussionActiveUser(c context.Context, limit int) ([]response.ActiveUserItem, error) {
	key := consts.GetGithubActiveUsersCacheKey(limit)
	var cached []response.ActiveUserItem
	ok, err := s.cached.GetJSON(c, key, &cached)
	if err != nil {
		return nil, err

	}
	if ok {
		return cached, nil
	}
	res, err := s.ds.GetActiveUser(c, limit)
	if err != nil {
		return nil, err
	}
	_ = s.cached.SetJSON(c, key, res, consts.CacheExpireDuration)
	return res, nil
}
