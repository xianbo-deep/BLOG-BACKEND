package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"context"
)

func GetSlowPages(ctx context.Context, limit int) ([]response.SlowDelayItem, error) {
	return dao.GetSlowPages(ctx, limit)
}

func GetAverageDelay() ([]response.AverageDelayItem, error) {
	return dao.GetAverageDelay()
}
