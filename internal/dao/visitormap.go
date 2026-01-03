package dao

import (
	"Blog-Backend/consts"
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

func GetChineseVisitorMap(startTime, endTime *time.Time) ([]response.ChineseVisitorMapItem, error) {
	var results []response.ChineseVisitorMapItem

	db := core.DB.Model(&model.VisitLog{})
	if startTime != nil {
		db = db.Where("visit_time >= ?", *startTime)
	}
	if endTime != nil {
		db = db.Where("visit_time <= ?", *endTime)
	}

	err := db.Select("region, count(*) as visitors").
		Where("country = ?", consts.CountryChina). // TODO 改成china
		Group("region").
		Scan(&results).Error

	return results, err
}
