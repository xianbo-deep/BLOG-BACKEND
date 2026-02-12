package discReport

import (
	"Blog-Backend/bootstrap"
	"Blog-Backend/consts"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/thirdparty/github/service"
	"context"
	"log"
	"time"
)

type DiscussionDigest struct {
	mailer *email.Mailer
	svc    *service.DiscussionService
}

func NewDiscussionDigest(cmp *bootstrap.Components) *DiscussionDigest {
	return &DiscussionDigest{mailer: cmp.Mailer, svc: cmp.GithubSVC}
}

func (d *DiscussionDigest) Start() {
	// 获取起始时间和结束时间
	startAt := consts.TransferTimeByLoc(time.Now().AddDate(0, 0, -7))
	endAt := consts.TransferTimeByLoc(time.Now())
	// 超时
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	digest, err := d.svc.GetDiscussionDigest(ctx, startAt, endAt)
	if err != nil {
		log.Println("获取评论区周报信息失败: %v", err)
		return
	}

	if err = d.mailer.SendTemplate([]string{consts.MyTencentEmail}, email.MailDiscussionDigest, digest, true); err != nil {
		log.Printf("发送评论区周报邮件失败: %v", err)
		return
	}
}
