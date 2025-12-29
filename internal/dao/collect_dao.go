package dao

import (
	"Blog-Backend/core"
	"Blog-Backend/model"
	"errors"
)

func InsertVisitLog(log model.VisitLog) error {
	/* 插入数据 */
	if core.DB != nil {
		core.DB.Create(&log)
		return nil
	}
	return errors.New("InsertVisitLog failed")
}

func IncrementPV(path string) error {
	/* 插入数据 */
	if core.RDB != nil {

		return nil
	}
	return errors.New("IncrementPV failed")
}

func IncrementUV(path string, visitorID string) error {
	/* 插入数据 */
	if core.RDB != nil {

		return nil
	}
	return errors.New("IncrementUV failed")
}
