package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/thirdparty/github/service"
	"context"
	"errors"
	"log"
	"time"
)

type GithubWebhookService struct {
	mailer    *email.Mailer
	githubsvc *service.DiscussionService
	dao       *dao.GithubWebhookDao
}

func NewGithubWebhookService(githubsvc *service.DiscussionService, dao *dao.GithubWebhookDao) *GithubWebhookService {
	mailer := email.RegisterEmail()
	return &GithubWebhookService{mailer: mailer, githubsvc: githubsvc, dao: dao}
}

func (s *GithubWebhookService) GetNewNotify(c context.Context) error {
	// 获取新的评论信息
	res, err := s.githubsvc.GetNewFeed(c, 1)
	if err != nil {
		log.Printf("无法获取最新评论信息: %v", err)
		return err
	}
	// 组装成模板所需要的
	data, ok := processData(res)
	if !ok {
		return errors.New("组装模板失败")
	}
	// 发送邮件通知
	err = s.mailer.SendTemplate([]string{consts.MyTencentEmail}, email.MailDiscussionNotify, data)
	if err != nil {
		log.Printf("最新评论邮件通知发送失败: %v", err)
		return err
	}
	return nil
}

func processData(items []*response.NewFeedItem) (email.DiscussionNotify, bool) {
	if len(items) == 0 || items[0] == nil {
		return email.DiscussionNotify{}, false
	}
	notify := items[0]
	res := email.DiscussionNotify{
		Type:           notify.EventType,
		User:           notify.Name,
		DiscussionTime: notify.Time,
		Text:           notify.Content,
		Avatar:         notify.Avatar,
		PageURL:        notify.Path,
		ReplyToAvatar:  notify.ReplyToAvatar,
		ReplyToUser:    notify.ReplyToName,
		ReplyToMessage: notify.ReplyToContent,
		FormattedTime:  consts.TransferTimeByLoc(notify.Time).Format(consts.TimeLayout),
		Year:           time.Now().Year(),
	}
	return res, true
}

func (s *GithubWebhookService) NotifySubscribeUsers(pages []email.ChangedPage, updatedAt time.Time, author string) error {
	res, err := s.dao.GetSubscribeUsers()
	if err != nil {
		return err
	}
	emails := make([]string, len(res))
	ids := make([]int64, len(res))
	for _, u := range res {
		emails = append(emails, u.Email)
		ids = append(ids, u.ID)
	}
	data := email.SubscribeNotify{
		Pages:               pages,
		UpdatedAt:           updatedAt,
		Author:              author,
		FormattedUpdateTime: updatedAt.Format(consts.TimeWithoutSecond),
		Year:                time.Now().Year(),
	}
	// 发送邮件
	e := s.mailer.SendTemplate(emails, email.MailDiscussionNotify, data)
	if e != nil {
		log.Printf("发送订阅邮件失败: %v", e)
		return e
	}
	// 更新订阅用户信息
	e = s.dao.UpdateSubscribeUsersLastSentTime(ids)
	if e != nil {
		log.Printf("更新用户信息失败: %v", e)
		return e
	}
	return nil
}
