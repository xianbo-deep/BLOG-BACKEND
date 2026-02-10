package email

import "time"

// 环境变量
const (
	EnvEmailHost = "EMAIL_HOST"
	EnvEmailPort = "EMAIL_PORT"
	EnvEmailSMTP = "EMAIL_SMTP"
	EnvEmailFrom = "EMAIL_FROM"
	EnvEmailUser = "EMAIL_USER"
)

// 邮件类型
const (
	MailDeadlinkReport   = "deadlink_report"
	MailDiscussionNotify = "discussion_notify"
	MailDiscussionDigest = "discussion_digest"
	MailSubscribeNotify  = "subscribe_notify"
	MailSubscribe        = "subscribe"
	MailUnSubscribe      = "unsubscribe"
)

// 标题
const (
	DeadLinkSubject         = "博客死链检测报告"
	DiscussionNotifySubject = "博客有新评论"
	DiscussionDigestSubject = "博客评论区周报"
	SubscribeNotifySubject  = "您订阅的博客更新了"
	SubscribeSubject        = "您已成功订阅"
	UnSubscribeSubject      = "您已取消订阅"
)

// 文件路径
const (
	DeadLinkFile         = "./template/deadlink_report.html"
	DiscussionReportFile = "./template/discussion_report.html"
	DiscussionNotifyFile = "./template/discussion_notify.html"
	SubscribeNotifyFile  = "./template/subscribe_notify.html"
	SubscribeFile        = "./template/subscribe.html"
	UnSubscribeFile      = "./template/unsubscribe.html"
)

// 页面改变类型
const (
	Added    = "Added"
	Modified = "Modified"
	Removed  = "Removed"
)

type EmailConfig struct {
	Host string // 执行发送的主机
	Port int    // 执行发送的端口
	User string // 发送的用户邮箱
	Pass string // SMTP密钥
	From string
}

// 死链检测
type DeadLinkReport struct {
	FinishedAt   time.Time
	PagesScanned int
	DeadlinkCnt  int
	LinksChecked int
	Items        []DeadLinkItem
}

type DeadLinkItem struct {
	FromPage   string
	LinkURL    string
	StatusCode int
	Err        string
}

// 评论通知
type DiscussionNotify struct {
	Type           string
	User           string
	DiscussionTime time.Time
	Avatar         string
	PageURL        string
	Text           string
	ReplyToUser    string
	ReplyToMessage string
	ReplyToAvatar  string
}

// 评论汇总
type DiscussionDigest struct {
	StartTime     time.Time
	EndTime       time.Time
	Comments      int
	Reactions     int
	Replies       int
	CommentItems  []CommentItem
	ReplyItems    []ReplyItem
	ReactionItems []ReactionItem
}

type CommentItem struct {
	User        string
	Avatar      string
	CommentTime time.Time
	PageURL     string
	Text        string
}

type ReplyItem struct {
	User           string
	Avatar         string
	ReplyTime      time.Time
	Text           string
	ReplyToUser    string
	ReplyToAvatar  string
	ReplyToMessage string
	PageURL        string
}

type ReactionItem struct {
	User         string
	Avatar       string
	ReactionTime time.Time
	PageURL      string
	ReactionType string
}

// 订阅通知
type SubscribeNotify struct {
	Pages     []ChangedPage
	UpdatedAt time.Time
	Author    string
}

type ChangedPage struct {
	Page       string
	ChangeType string
	Path       string
}

// 订阅与取消订阅
type SubscribeOrNot struct {
	Year int
}
