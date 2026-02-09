package public

import (
	"Blog-Backend/internal/dao"
	"log"
)

type SubscribeService struct {
	dao *dao.SubscribeDao
}

func NewSubscribeService(dao *dao.SubscribeDao) *SubscribeService {
	return &SubscribeService{dao: dao}
}

func (s *SubscribeService) SubscribeBlog(email string, subscribe int) error {
	err := s.dao.SubscribeBlog(email, subscribe)
	if err != nil {
		log.Printf("订阅邮件dao层出现错误: %v", err)
		return err
	}
	return nil
}
