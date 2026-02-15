package deadlink

import "time"

const (
	defaultConcurrency = 10
	retryTimeout       = 1 * time.Second
	failedMsg          = "死链检测失败"
	sitemapURLSuffix   = "/sitemap.xml"
	retryTimes         = 5
	timeout            = 5 * time.Second
	cloneRetryTimes    = 3

	RepoURL = "https://github.com/xianbo-deep/xbZhong.git"
	Branch  = "main"
	DocsDir = "docs"

	CacheRepoDir = "/var/cache/deadlink/repo.git"
	ProxyHTTP    = "http://127.0.0.1:7890"
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

	RepoURL string // 仓库地址
	Branch  string // 分支
	DocsDir string // 根目录

	CacheRepoDir string // 仓库缓存目录
	ProxyHTTP    string // 代理
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
