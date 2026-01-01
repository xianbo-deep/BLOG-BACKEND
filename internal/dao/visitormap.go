package dao

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"time"
)

func GetVisitorMap(startTime, endTime *time.Time) ([]response.VisitorMapItem, error) {
	var results []response.VisitorMapItem

	db := core.DB.Model(&model.VisitLog{})

	if startTime != nil {
		db = db.Where("visit_time >= ?", *startTime)
	}
	if endTime != nil {
		db = db.Where("visit_time <= ?", *endTime)
	}

	err := db.Select("country, count(*) as visitors").
		Group("country").
		Scan(&results).Error

	return results, err
}
