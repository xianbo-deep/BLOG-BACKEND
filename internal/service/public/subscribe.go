package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/notify/email"
	"errors"
	"log"
	"time"
)

type SubscribeService struct {
	dao    *dao.SubscribeDao
	mailer *email.Mailer
}

func NewSubscribeService(dao *dao.SubscribeDao, mailer *email.Mailer) *SubscribeService {
	return &SubscribeService{dao: dao, mailer: mailer}
}

func (s *SubscribeService) SubscribeBlog(mail string, subscribe int) error {
	// 存储订阅信息
	err := s.dao.SubscribeBlog(mail, subscribe)
	if err != nil {
		log.Printf("订阅邮件dao层出现错误: %v", err)
		return err
	}

	data := email.SubscribeOrNot{Year: consts.TransferTimeByLoc(time.Now()).Year()}
	// 发送邮件通知
	if subscribe == consts.Subsucribed {
		e := s.mailer.SendTemplate([]string{mail}, email.MailSubscribe, data)
		return e
	} else if subscribe == consts.UnSubsucribed {
		e := s.mailer.SendTemplate([]string{mail}, email.MailUnSubscribe, data)
		return e
	}
	return errors.New("找不到订阅类型")
}
