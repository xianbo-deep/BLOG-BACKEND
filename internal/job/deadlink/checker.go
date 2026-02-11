package deadlink

import (
	"Blog-Backend/consts"
	"Blog-Backend/internal/notify/email"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Checker struct {
	cfg    Config
	client *http.Client
	mailer *email.Mailer
}

func NewChecker(cfg Config, mailer *email.Mailer) *Checker {
	if cfg.Concurrency <= 0 || cfg.Concurrency > 100 {
		cfg.Concurrency = 10
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = time.Second * 5
	}
	return &Checker{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
		mailer: mailer,
	}
}

func (c *Checker) Check() (Summary, []Result, error) {
	sum := Summary{
		StartedAT: time.Now().UTC(),
	}

	// 得到页面列表
	pages, err := c.fetchSitemapURLs(c.cfg.SitemapURL)
	if err != nil {
		return Summary{}, nil, err
	}

	// 得到链接列表
	links, pagesScanned, err := c.collectLinksFromPages(pages)
	if err != nil {
		return Summary{}, nil, err
	}

	// 获取要扫描的页面数
	sum.PagesScanned = pagesScanned

	// 提取总检测数
	sum.LinksChecked = len(links)

	// 获取检测结果
	results := c.checkLinks(links)

	// 获取死链数量
	for _, res := range results {
		if !res.OK {
			sum.DeadlinkCnt++
		}
	}

	// 更新结束时间
	sum.FinishedAT = time.Now().UTC()

	return sum, results, nil
}

func (c *Checker) fetchSitemapURLs(sitemapurl string) ([]string, error) {
	resp, err := c.client.Get(sitemapurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var set SitemapURLSet
	// 解析
	if err := xml.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, err
	}
	urls := make([]string, 0, len(set.URLs))
	for _, u := range set.URLs {
		loc := strings.TrimSpace(u.Loc)
		urls = append(urls, loc)
	}
	return urls, nil
}

// 从页面获取链接
func (c *Checker) collectLinksFromPages(pages []string) ([]LinkPair, int, error) {
	// 协程数
	n := c.cfg.Concurrency

	// 创建通道用于传输页面
	jobs := make(chan string, 2*n)

	// 创建link的map，用于判断是否检测
	collected := make(map[string]struct{}, 4096)

	// 创建任务组
	var wg sync.WaitGroup

	// 创建锁
	var mu sync.Mutex

	// 进行检测的页面数
	var pagesScanned int32

	// 创建返回值
	res := make([]LinkPair, 0, 2048)

	worker := func() {
		defer wg.Done()
		for page := range jobs {
			links, err := c.extractLinksFromHTML(page)
			if err != nil {
				// TODO 抓取某个页面的链接失败
				continue
			}
			atomic.AddInt32(&pagesScanned, 1)
			for _, link := range links {
				key := link
				// 加锁防止并发操作map导致panic
				mu.Lock()
				if _, ok := collected[key]; ok {
					mu.Unlock()
					continue
				}
				collected[key] = struct{}{}
				res = append(res, LinkPair{fromPage: page, linkURL: link})
				mu.Unlock()

			}
		}
	}

	// 添加协程

	wg.Add(n)

	for i := 0; i < n; i++ {
		go worker()
	}

	for _, page := range pages {
		jobs <- page
	}

	close(jobs)
	// 阻塞直到任务完成
	wg.Wait()

	return res, int(pagesScanned), nil
}

// 从HTML中抽取links
func (c *Checker) extractLinksFromHTML(pageURL string) ([]string, error) {
	// 创建请求体
	req, _ := http.NewRequest("GET", pageURL, nil)
	req.Header.Set("User-Agent", "DeadlinkChecker/1.0(+xbzhong.cn)")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// 获取得到的html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	seen := map[string]struct{}{}
	out := make([]string, 0, 64)

	// 抓取链接
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		// 提取href标签
		href, _ := s.Attr("href")
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}
		// 转小写
		lower := strings.ToLower(href)

		// 过滤无效链接
		if strings.HasPrefix(lower, "#") ||
			strings.HasPrefix(lower, "javascript:") ||
			strings.HasPrefix(lower, "mailto:") ||
			strings.HasPrefix(lower, "tel:") {
			return
		}

		// 解析url
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		// 相对路径转成绝对路径
		absURL := base.ResolveReference(u)

		// 过滤站内链接
		if base.Host == absURL.Host {
			return
		}

		abs := absURL.String()
		// 去重
		if _, ok := seen[abs]; ok {
			return
		}

		seen[abs] = struct{}{}
		out = append(out, abs)
	})
	return out, nil
}

// 检测链接
func (c *Checker) checkLinks(links []LinkPair) []Result {
	n := c.cfg.Concurrency
	if n <= 0 {
		n = defaultConcurrency
	}

	jobs := make(chan LinkPair, n*2)
	out := make(chan Result, 1024)

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		// 检查页面中的链接
		for p := range jobs {
			statusCode, ok, errStr := c.checkLink(p.linkURL)
			out <- Result{
				FromPage:   p.fromPage,
				LinkURL:    p.linkURL,
				OK:         ok,
				StatusCode: statusCode,
				Err:        errStr,
				CheckedAT:  time.Now().UTC(),
			}
		}
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}

	for _, link := range links {
		jobs <- link
	}
	close(jobs)

	// 开一个协程关闭out，让主协程立马消费out
	go func() {
		wg.Wait()
		close(out)
	}()

	res := make([]Result, 0, len(links))
	for r := range out {
		res = append(res, r)
	}

	return res
}

func (c *Checker) checkLink(link string) (status int, ok bool, errStr string) {
	for attempt := 0; attempt <= c.cfg.Retry; attempt++ {
		status, ok, errStr := c.headThenGet(link)
		if ok {
			return status, true, ""
		}
		if status >= 400 && status < 500 && status != http.StatusRequestTimeout {
			return status, false, errStr
		}
		if attempt < c.cfg.Retry {
			// 防止频繁爬导致被封IP
			time.Sleep(retryTimeout)
		}
	}
	return status, false, errStr
}

func (c *Checker) headThenGet(link string) (status int, ok bool, errStr string) {
	// Head
	req, _ := http.NewRequest("HEAD", link, nil)
	req.Header.Set("User-Agent", "DeadlinkChecker/1.0(+xbzhong.cn)")
	resp, err := c.client.Do(req)
	if err == nil && resp != nil {
		defer resp.Body.Close()
		status = resp.StatusCode
		if status >= 200 && status < 400 {
			return status, true, ""
		}
		if status != http.StatusMethodNotAllowed && status != http.StatusForbidden {
			return status, false, ""
		}
	} else if err != nil {
		return status, false, err.Error()
	}
	// Get
	req, _ = http.NewRequest("GET", link, nil)
	req.Header.Set("User-Agent", "DeadlinkChecker/1.0(+xbzhong.cn)")
	resp, err = c.client.Do(req)
	if err != nil {
		return 0, false, err.Error()
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	if status >= 200 && status < 400 {
		return status, true, ""
	}
	return status, false, failedMsg
}

// 组装成模板需要的结构体
func (c *Checker) processData(summary Summary, results []Result) DeadLinkReportData {
	var data DeadLinkReportData
	// 组装全局信息
	data.PagesScanned = summary.PagesScanned
	data.LinksChecked = summary.LinksChecked
	data.BJTime = consts.TransferTimeByLoc(summary.FinishedAT).Format(consts.TimeWithoutSecond)
	data.Year = consts.TransferTimeByLoc(time.Now()).Year()

	// 组装详细信息
	var deadLinks []DeadLinkItem
	for _, item := range results {
		if item.OK {
			continue
		}
		deadLinks = append(deadLinks, DeadLinkItem{
			Page:    item.FromPage,
			URL:     item.LinkURL,
			Status:  item.StatusCode,
			Message: item.Err,
		})
	}

	return data
}
