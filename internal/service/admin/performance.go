package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"context"
)

type PerformanceService struct {
	dao *dao.PerformanceDao
}

func NewPerformanceService(dao *dao.PerformanceDao) *PerformanceService {
	return &PerformanceService{dao: dao}
}

func (s *PerformanceService) GetSlowPages(ctx context.Context, limit int) ([]response.SlowDelayItem, error) {
	return s.dao.GetSlowPages(ctx, limit)
}

func (s *PerformanceService) GetAverageDelay() ([]response.AverageDelayItem, error) {
	return s.dao.GetAverageDelay()
}
