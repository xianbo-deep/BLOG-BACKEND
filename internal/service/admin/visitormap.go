package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"time"
)

type VisitorMapSerive struct {
	dao *dao.VisitorMapDao
}

func NewVisitorMapSerive(dao *dao.VisitorMapDao) *VisitorMapSerive {
	return &VisitorMapSerive{dao: dao}
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

	return s.dao.GetVisitorMap(start, end)
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

	return s.dao.GetChineseVisitorMap(start, end)

}
