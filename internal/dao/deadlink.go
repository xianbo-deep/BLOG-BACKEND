package dao

import (
	"Blog-Backend/model"

	"gorm.io/gorm"
)

type DeadLinkDao struct {
	db *gorm.DB
}

func NewDeadLinkDao(db *gorm.DB) *DeadLinkDao {
	return &DeadLinkDao{db: db}
}

func (d *DeadLinkDao) SaveRunAndItems(run model.DeadLinkRun, items []model.DeadLinkItem) error {
	var runID int64

	// 开启事务
	err := d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&run).Error; err != nil {
			return err
		}
		runID = run.ID

		// 无死链检测
		if len(items) == 0 {
			return nil
		}

		// 绑定执行ID
		for i := range items {
			items[i].RunID = runID
		}

		// 批量插入
		if err := tx.CreateInBatches(items, 2000).Error; err != nil {
			return err
		}
		return nil
	})

	return err

}
