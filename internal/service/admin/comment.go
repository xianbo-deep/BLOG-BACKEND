package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao/cache"
	"Blog-Backend/thirdparty/github/service"
	"context"
)

type CommentService struct {
	cached *cache.CacheDAO
	ds     *service.DiscussionService
}

func NewCommentService(cached *cache.CacheDAO, ds *service.DiscussionService) *CommentService {
	return &CommentService{cached: cached, ds: ds}
}

func (s *CommentService) GetDiscussionMetric(c context.Context, days int) (response.Metric, error) {
	version, err := s.cached.GetVersion(c, consts.RedisGithubCacheVerKey)
	if err != nil || version == -1 {
		return response.Metric{}, err
	}
	key := consts.GetGithubMetricCacheKey(version, days)
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
	version, err := s.cached.GetVersion(c, consts.RedisGithubCacheVerKey)
	if err != nil || version == -1 {
		return nil, err
	}
	key := consts.GetGithubTrendCacheKey(version, days)
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
	version, err := s.cached.GetVersion(c, consts.RedisGithubCacheVerKey)
	if err != nil || version == -1 {
		return nil, err
	}
	key := consts.GetGithubNewFeedsCacheKey(version, limit)
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
	version, err := s.cached.GetVersion(c, consts.RedisGithubCacheVerKey)
	if err != nil {
		return nil, err
	}
	key := consts.GetGithubActiveUsersCacheKey(version, limit)
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
