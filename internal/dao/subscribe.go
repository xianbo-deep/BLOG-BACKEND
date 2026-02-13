package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SubscribeDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewSubscribeDao(db *gorm.DB, rdb *redis.Client) *SubscribeDao {
	return &SubscribeDao{db: db, rdb: rdb}
}

func (d *SubscribeDao) SubscribeBlog(email string, subscribe int) error {
	var user model.SubscribeUser
	err := d.db.Where("email = ?", email).First(&user).Error
	// 用户已经存在
	if err == nil {
		if user.Status == subscribe {
			return errors.New("之前已经成功/取消订阅，请勿重复操作")
		}
		return d.db.Model(&model.SubscribeUser{}).
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

func (d *SubscribeDao) StoreVC(ctx context.Context, email, vc string) error {
	// 获取缓存key
	key := consts.VerificationCodeKey + email
	// 在redis中存入验证码
	if err := d.rdb.Set(ctx, key, vc, 5*consts.TimeRangeMinute).Err(); err != nil {
		return fmt.Errorf("存储验证码到redis失败: %w", err)
	}
	return nil
}

func (d *SubscribeDao) VerifyVC(ctx context.Context, email, vc string) error {
	// 获取缓存key
	key := consts.VerificationCodeKey + email

	// 从redis获取验证码
	storedVC, err := d.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("验证码过期或不存在")
	} else if err != nil {
		return err
	}

	if storedVC != vc {
		return fmt.Errorf("验证码正确")
	}

	// 验证成功
	return nil
}

func (d *SubscribeDao) DelVC(ctx context.Context, email string) error {
	key := consts.VerificationCodeKey + email
	// 删除验证码
	if err := d.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
