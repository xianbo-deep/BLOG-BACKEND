package discReport

import (
	"Blog-Backend/consts"
	"Blog-Backend/internal/notify/email"
	"Blog-Backend/thirdparty/github/service"
	"time"
)

type DiscussionDigest struct {
	mailer *email.Mailer
	svc    *service.DiscussionService
}

func NewDiscussionDigest(mailer *email.Mailer, svc *service.DiscussionService) *DiscussionDigest {
	return &DiscussionDigest{mailer: mailer, svc: svc}
}

func (d *DiscussionDigest) Start() {
	startAt := consts.TransferTimeByLoc(time.Now().AddDate(0, 0, -7))
	endAt := consts.TransferTimeByLoc(time.Now())
	digest := d.GetDiscussionDigest(startAt, endAt)
}
func (d *DiscussionDigest) GetDiscussionDigest(startAt, endAt time.Time) email.DiscussionDigest {

}

func (d *DiscussionDigest) SendEmail() {

}
