package deadlink

import "time"

const (
	defaultConcurrency = 10
	retryTimeout       = 200 * time.Millisecond
	failedMsg          = "死链检测失败"
	sitemapURLSuffix   = "/sitemap.xml"
	retryTimes         = 5
	timeout            = 5 * time.Second
)

type Summary struct {
	StartedAT    time.Time
	FinishedAT   time.Time
	PagesScanned int
	DeadlinkCnt  int
	LinksChecked int
}

type Result struct {
	FromPage   string
	LinkURL    string
	StatusCode int
	OK         bool
	Err        string
	CheckedAT  time.Time
}

// 死链检测配置
type Config struct {
	SitemapURL  string        // 博客sitemap链接
	Concurrency int           // 并发数量
	Timeout     time.Duration // 超时时间
	Retry       int           // 重试次数
	URLPrefix   string        // 网页前缀
}

// sitemap标准格式
// vuepress的标准格式
type SitemapURLSet struct {
	URLs []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

type LinkPair struct {
	fromPage string
	linkURL  string
}

// 完善模板结构体
type DeadLinkReportData struct {
	BJTime       string
	Year         int
	PagesScanned int
	DeadLinkCnt  int
	LinksChecked int
	DeadLinks    []DeadLinkItem
}

type DeadLinkItem struct {
	Page    string
	Status  int
	URL     string
	Message string
}
