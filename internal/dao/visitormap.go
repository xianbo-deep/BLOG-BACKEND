package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"time"

	"gorm.io/gorm"
)

type VisitorMapDao struct {
	db *gorm.DB
}

func NewVisitorMapDao(db *gorm.DB) *VisitorMapDao {
	return &VisitorMapDao{db: db}
}

func (d *VisitorMapDao) GetVisitorMap(startTime, endTime *time.Time) ([]response.VisitorMapItem, error) {
	var results []response.VisitorMapItem

	db := d.db.Model(&model.VisitLog{})

	if startTime != nil {
		db = db.Where("visit_time >= ?", *startTime)
	}
	if endTime != nil {
		db = db.Where("visit_time <= ?", *endTime)
	}

	err := db.Select("country_en as country, count(*) as visitors").
		Group("country_en").
		Scan(&results).Error

	return results, err
}

func (d *VisitorMapDao) GetChineseVisitorMap(startTime, endTime *time.Time) ([]response.ChineseVisitorMapItem, error) {
	var results []response.ChineseVisitorMapItem

	db := d.db.Model(&model.VisitLog{})
	if startTime != nil {
		db = db.Where("visit_time >= ?", *startTime)
	}
	if endTime != nil {
		db = db.Where("visit_time <= ?", *endTime)
	}

	err := db.Select("region, count(*) as visitors").
		Where("country_code = ?", consts.CountryChinaCode).
		Group("region").
		Scan(&results).Error

	return results, err
}
