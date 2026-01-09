package response

type Metric struct {
	TotalComments  int64 `json:"total_comments"`  // 总评论数
	TotalReplies   int64 `json:"total_replies"`   // 总回复数
	TotalReactions int64 `json:"total_reactions"` // 总回应数
}

type NewFeedItem struct {
	Path           string `json:"path"`           // 动态页面路径
	URL            string `json:"url"`            // 用户主页
	Name           string `json:"name"`           // 用户名称
	EventType      string `json:"event_type"`     // 事件类型
	Avatar         string `json:"avatar"`         // 用户头像
	Time           string `json:"time"`           // 动态发生时间
	Content        string `json:"content"`        // 动态内容
	ReplyToName    string `json:"replyToName"`    // 被回复者的名字
	ReplyToAvatar  string `json:"replyToAvatar"`  // 被回复者头像
	ReplyToContent string `json:"replyToContent"` // 被回复者评论内容
}

type ActiveUserItem struct {
	Name       string `json:"name"`       // 用户名字
	TotalFeeds int64  `json:"totalFeeds"` // 用户贡献动态数
	Avatar     string `json:"avatar"`     // 用户头像
	URL        string `json:"URL"`        // 用户主页URL
}

type TrendItem struct {
	Date           string `json:"date"`           // 日期
	TotalComments  int64  `json:"totalComments"`  // 总评论数
	TotalReplies   int64  `json:"totalReplies"`   // 总回复数
	TotalReactions int64  `json:"totalResponses"` // 总回应数
}
