package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/notify/email"
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type SubscribeService struct {
	dao    *dao.SubscribeDao
	mailer *email.Mailer
}

func NewSubscribeService(dao *dao.SubscribeDao, mailer *email.Mailer) *SubscribeService {
	return &SubscribeService{dao: dao, mailer: mailer}
}

func (s *SubscribeService) SubscribeBlog(ctx context.Context, mail, vc string, subscribe int) error {
	// 获取超时上下文
	c, cancel := context.WithTimeout(ctx, 2*consts.TimeRangeSecond)
	defer cancel()
	// 验证验证码
	e := s.dao.VerifyVC(c, mail, vc)
	if e != nil {
		return e
	}

	// 存储订阅信息
	err := s.dao.SubscribeBlog(mail, subscribe)
	if err != nil {
		log.Printf("订阅邮件dao层出现错误: %v", err)
		return err
	}

	// 删除验证码
	err = s.dao.DelVC(c, vc)
	if err != nil {
		log.Printf("验证码删除失败: %v", err)
	}
	data := email.SubscribeOrNot{Year: consts.TransferTimeByLoc(time.Now()).Year()}
	// 发送邮件通知
	if subscribe == consts.Subscribed {
		e := s.mailer.SendTemplate([]string{mail}, email.MailSubscribe, data, true)
		return e
	} else if subscribe == consts.UnSubscribed {
		e := s.mailer.SendTemplate([]string{mail}, email.MailUnSubscribe, data, true)
		return e
	}

	return errors.New("找不到订阅类型")
}

func (s *SubscribeService) VerifyEmail(ctx context.Context, mail string, subscribe int) error {
	// 获取验证码
	vc := s.generateVC()
	// 存储验证码
	err := s.dao.StoreVC(ctx, mail, vc)
	if err != nil {
		log.Printf("存储验证码到redis失败: %v", err)
		return err
	}
	// 准备模板结构体
	data := email.SubscribeVerificationCode{
		Year:      consts.TransferTimeByLoc(time.Now()).Year(),
		VC:        vc,
		Email:     mail,
		Subscribe: (subscribe == consts.Subscribed),
	}
	// 发送邮箱
	err = s.mailer.SendTemplate([]string{mail}, email.MailSubscribeVerify, data, true)
	if err != nil {
		log.Printf("发送验证码邮件失败: %v", err)
		return err
	}
	return nil
}

func (s *SubscribeService) generateVC() string {
	rand.Seed(time.Now().UnixNano())
	vc := rand.Intn(900000) + 100000
	return strconv.Itoa(vc)
}
