package discReport

import (
	"Blog-Backend/bootstrap"

	"github.com/robfig/cron/v3"
)

func RegisterDiscussionDigest(c *cron.Cron, cmp *bootstrap.Components) {
	c.AddFunc("0 0 8 * * MON", func() {

	})
}
