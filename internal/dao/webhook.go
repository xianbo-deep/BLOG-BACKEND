package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GithubWebhookDao struct {
	db *gorm.DB
}

func NewGithubWebhookDao(db *gorm.DB) *GithubWebhookDao {
	return &GithubWebhookDao{db: db}
}

func (d *GithubWebhookDao) GetSubscribeUsers() ([]model.SubscribeUser, error) {
	var res []model.SubscribeUser

	db := d.db.Model(&model.SubscribeUser{})
	// 获取已订阅的用户
	err := db.Where("status = ?", consts.Subscribed).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 更新用户数据
func (d *GithubWebhookDao) UpdateSubscribeUsersLastSentTime(ids []int64) error {
	if len(ids) == 0 {
		return errors.New("订阅用户为空")
	}
	db := d.db.Model(&model.SubscribeUser{})

	return db.
		Where("id in (?)", ids).
		Update("last_sent_at", time.Now().UTC()).Error
}
