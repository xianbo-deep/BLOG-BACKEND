package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/model"

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
	err := db.Where("status = ?", consts.Subsucribed).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
