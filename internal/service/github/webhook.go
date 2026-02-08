package github

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/thirdparty/github/service"
	"context"
	"log"
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

func (s *GithubWebhookService) GetNewNotify(c context.Context) {
	// 获取新的评论信息
	res, err := s.githubsvc.GetNewFeed(c, 1)
	if err != nil {
		log.Printf("无法获取最新评论信息: %v", err)
		return
	}
	// 组装成模板所需要的
	data, ok := s.processData(res)
	if !ok {
		return
	}
	// 发送邮件通知
	err = s.mailer.SendTemplate([]string{consts.MyTencentEmail}, email.MailDiscussionNotify, data)
	if err != nil {
		log.Printf("最新评论邮件通知发送失败: %v", err)
		return
	}
}

func (s *GithubWebhookService) processData(items []*response.NewFeedItem) (email.DiscussionNotify, bool) {
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
	}
	return res, true
}

func (s *GithubWebhookService) NotifySubscribeUser() {

}
