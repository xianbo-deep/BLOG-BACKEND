package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
)

func GetSlowPages(limit int) ([]response.SlowDelayItem, error) {
	res, err := dao.GetSlowPages(limit)
	if err != nil {
		return []response.SlowDelayItem{}, err
	}

	return res, nil
}

func GetAverageDelay() ([]response.AverageDelayItem, error) {
	res, err := dao.GetAverageDelay()
	if err != nil {
		return []response.AverageDelayItem{}, err
	}

	return res, nil
}
