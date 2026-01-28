package deadlink

import (
	"net/http"
	"time"
)

// 死链检测配置
type Config struct {
	SitemapURL  string
	Concurrency int
	Timeout     time.Duration
	Retry       int
}

// 结果
type Result struct {
	FromPage   string
	CheckURL   string
	StatusCode int
	OK         bool
	Err        string
	CheckedAT  time.Time
}

// 汇总
type Summary struct {
	StartedAT    time.Time
	FinishedAT   time.Time
	PagesScanned int
	LinksChecked int
	DeadCount    int
}

type Checker struct {
	cfg    Config
	client *http.Client
}

func NewChecker(cfg Config) *Checker {
	if cfg.Concurrency < 0 {
		cfg.Concurrency = 10
	}
	if cfg.Timeout < 0 {
		cfg.Timeout = 10 * time.Second
	}
	return &Checker{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// sitemapURL列表
type sitemapURLSet struct {
	URLs[] struct{
		Loc string `xml:"loc"`
	} `xml:"url"`
}

func (c *Checker) Run() (Summary,[]Result,error){
	// 初始化起始时间
	sum := Summary{
		StartedAT: time.Now(),
	}

	runerr := error(nil)

	// 拿sitemap的所有文章
	c.client.
	// 遍历文章里面的链接

	// 写入返回值

}

func (c *Checker) fetchSitemapPages(sitemapURL string) ([]sitemapURLSet, error) {
	resp,err := c.client.Get(sitemapURL)
	if err != nil {
		return nil, err
	}
	// 关闭
	defer resp.Body.Close()

	
}
