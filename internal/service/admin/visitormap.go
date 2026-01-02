package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"time"
)

type VisitorMapSerive struct {
}

func NewVisitorMapSerive() *VisitorMapSerive {
	return &VisitorMapSerive{}
}

func (s *VisitorMapSerive) GetVisitorMap(startTime int, endTime int) ([]response.VisitorMapItem, error) {
	var start, end *time.Time
	// 将时间戳转换为 time.Time
	if startTime != 0 {
		s := time.UnixMilli(int64(startTime))
		start = &s
	}
	if endTime != 0 {
		e := time.UnixMilli(int64(endTime))
		end = &e
	}

	return dao.GetVisitorMap(start, end)
}

func (s *VisitorMapSerive) GetChineseVisitorMap(startTime int, endTime int) ([]response.ChineseVisitorMapItem, error) {
	var start, end *time.Time

	if startTime != 0 {
		s := time.UnixMilli(int64(startTime))
		start = &s
	}
	if endTime != 0 {
		e := time.UnixMilli(int64(endTime))
		end = &e
	}

	return dao.GetChineseVisitorMap(start, end)

}
