package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"context"
)

type PerformanceService struct{}

func NewPerformanceService() *PerformanceService {
	return &PerformanceService{}
}

func (s *PerformanceService) GetSlowPages(ctx context.Context, limit int) ([]response.SlowDelayItem, error) {
	return dao.GetSlowPages(ctx, limit)
}

func (s *PerformanceService) GetAverageDelay() ([]response.AverageDelayItem, error) {
	return dao.GetAverageDelay()
}
