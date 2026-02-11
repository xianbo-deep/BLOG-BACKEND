package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SubscribeDao struct {
	db *gorm.DB
}

func NewSubscribeDao(db *gorm.DB) *SubscribeDao {
	return &SubscribeDao{db: db}
}

func (d *SubscribeDao) SubscribeBlog(email string, subscribe int) error {

	var user model.SubscribeUser
	err := d.db.Where("email = ?", email).First(&user).Error
	// 用户已经存在
	if err == nil {
		if user.Status == subscribe {
			return errors.New("之前已经成功/取消订阅，请勿重复操作")
		}
		return d.db.
			Where("id = ?", user.ID).
			Updates(map[string]any{
				"status":     uint8(subscribe),
				"updated_at": time.Now().UTC(),
			}).Error
	}
	// 其他错误直接返回
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if subscribe == consts.UnSubscribed {
		return errors.New("还未订阅，无法进行取消")
	}

	// 不存在则创建新记录
	newUser := model.SubscribeUser{
		Email:         email,
		Status:        subscribe,
		SubscribeTime: time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		NotifyCount:   0,
	}
	return d.db.Create(&newUser).Error
}
