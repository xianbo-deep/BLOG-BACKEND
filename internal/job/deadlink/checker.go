package deadlink

import (
	"Blog-Backend/consts"
	"Blog-Backend/internal/notify/email"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
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

	// 获取克隆仓库的路径
	log.Printf("开始克隆仓库")
	repoDir, err := c.cloneRepoToTemp()
	if err != nil {
		return Summary{}, nil, err
	}
	// 清理磁盘临时目录
	defer os.RemoveAll(repoDir)

	// 获取存储md文件的目录
	docsPath := filepath.Join(repoDir, c.cfg.DocsDir)

	// 得到链接列表
	log.Printf("获取外链列表")
	links, pagesScanned, err := c.collectLinksFromMarkdownDir(docsPath)
	if err != nil {
		return Summary{}, nil, err
	}

	// 获取要扫描的页面数
	sum.PagesScanned = pagesScanned

	// 提取总检测数
	sum.LinksChecked = len(links)

	// 获取检测结果
	log.Printf("检测外链")
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

// 克隆仓库到磁盘
func (c *Checker) cloneRepoToTemp() (string, error) {
	var lastErr error

	for i := 0; i < cloneRetryTimes; i++ {
		// 超时上下文
		ctx, cancel := consts.GetTimeoutContext(context.Background(), 3*time.Minute)

		// 创建临时目录
		dir, err := os.MkdirTemp("", "deadlink-repo-*")
		if err != nil {
			return "", err
		}
		// 执行克隆
		cmd := exec.CommandContext(
			ctx,
			"git", "clone",
			"--depth", "1",
			"--single-branch",
			"--branch", c.cfg.Branch,
			c.cfg.RepoURL,
			dir,
		)
		// 继承环境变量
		cmd.Env = os.Environ()

		// 记录输出
		out, err := cmd.CombinedOutput()

		// 关闭上下文
		cancel()

		if ctx.Err() == context.DeadlineExceeded {
			lastErr = fmt.Errorf("克隆超时: %s", string(out))
		} else if err != nil {
			lastErr = fmt.Errorf("克隆出现错误: %v, out = %s", err, string(out))
		} else {
			log.Printf("仓库克隆成功")
			return dir, nil
		}
		log.Printf("[deadlink] clone attempt=%d err=%v", i+1, lastErr)

		lastErr = err
		_ = os.RemoveAll(dir)
		time.Sleep(retryTimeout)
	}

	return "", lastErr

}

// 从md文件收集外部链接
func (c *Checker) collectLinksFromMarkdownDir(docsPath string) ([]LinkPair, int, error) {
	// 获取绝对路径
	absDocs, err := filepath.Abs(docsPath)
	if err != nil {
		return nil, 0, err
	}

	// 映射map
	collected := make(map[string]struct{}, 4096)

	// 输出结果
	out := make([]LinkPair, 0, 2048)

	// 扫描页面数
	pagesScanned := 0

	// 过滤文件(相对路径)
	skipFiles := map[string]struct{}{
		"README.md": {},
	}

	err = filepath.WalkDir(absDocs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			name := d.Name()
			if name == ".vuepress" || name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))

		// 判断是不是md文件
		if ext != ".md" && ext != ".mdx" {
			return nil
		}

		// 获取相对路径
		rel, _ := filepath.Rel(absDocs, path)
		rel = filepath.ToSlash(rel)

		// 过滤特定文件
		if _, ok := skipFiles[rel]; ok {
			return nil
		}

		// 获取文件内容
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// 增加检测文件数
		pagesScanned++

		links := c.extractLinksFromMarkdown(string(b))
		if len(links) == 0 {
			return nil
		}

		for _, link := range links {
			key := rel + "||" + link
			if _, ok := collected[key]; ok {
				continue
			}
			collected[key] = struct{}{}
			out = append(out, LinkPair{fromPage: rel, linkURL: link})
		}
		return nil
	})
	if err != nil {
		return nil, pagesScanned, err
	}
	return out, pagesScanned, nil
}

// 从文件中获取链接
func (c *Checker) extractLinksFromMarkdown(content string) []string {
	// 删除代码块
	reCode := regexp.MustCompile("(?s)```.*?```")
	content = reCode.ReplaceAllString(content, " ")

	// 删除行内代码
	reInlineCode := regexp.MustCompile("`[^`]*`")
	content = reInlineCode.ReplaceAllString(content, " ")

	// 删除图片
	reImg := regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)
	content = reImg.ReplaceAllString(content, " ")

	// 提取外链
	reURL := regexp.MustCompile(`https?://[A-Za-z0-9]+(?:\.[A-Za-z0-9]+)+`)
	raw := reURL.FindAllString(content, -1)

	if len(raw) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(raw))
	out := make([]string, 0, len(raw))

	for _, u := range raw {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		// 删除右边符号
		u = strings.TrimRightFunc(u, func(r rune) bool {
			// 结尾是标点就删（含中文标点）
			return unicode.IsPunct(r)
		})

		pu, err := url.Parse(u)
		if err != nil {
			continue
		}

		// 协议和主机名不能为空
		if pu.Scheme == "" || pu.Host == "" {
			continue
		}

		low := strings.ToLower(u)
		if strings.HasPrefix(low, "#") ||
			strings.HasPrefix(low, "mailto:") ||
			strings.HasPrefix(low, "tel:") ||
			strings.HasPrefix(low, "javascript:") {
			continue
		}

		if _, ok := seen[u]; ok {
			continue
		}

		seen[u] = struct{}{}
		out = append(out, u)
	}
	return out
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
		status, ok, errStr = c.headThenGet(link)
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
		status = resp.StatusCode
		drainAndClose(resp)
		if status >= 200 && status < 400 {
			return status, true, ""
		}
	} else if err != nil {
		errStr = err.Error()
	}
	// Get
	req, _ = http.NewRequest("GET", link, nil)
	req.Header.Set("User-Agent", "DeadlinkChecker/1.0(+xbzhong.cn)")
	resp, err = c.client.Do(req)
	if err != nil {
		return 0, false, err.Error()
	}
	status = resp.StatusCode
	drainAndClose(resp)
	if status >= 200 && status < 400 {
		return status, true, ""
	}
	if errStr != "" {
		return status, false, errStr
	}
	return status, false, failedMsg
}

// 把响应体读取到干净状态再关闭
func drainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20)) // 最多吞 1MB
	_ = resp.Body.Close()
}

// 组装成模板需要的结构体
func (c *Checker) processData(summary Summary, results []Result) email.DeadLinkReportData {
	var data email.DeadLinkReportData
	// 组装全局信息
	data.PagesScanned = summary.PagesScanned
	data.LinksChecked = summary.LinksChecked
	data.BJTime = consts.TransferTimeByLoc(summary.FinishedAT).Format(consts.TimeWithoutSecond)
	data.Year = consts.TransferTimeByLoc(time.Now()).Year()

	// 组装详细信息
	var deadLinks []email.DeadLinkItem
	for _, item := range results {
		if item.OK {
			continue
		}
		deadLinks = append(deadLinks, email.DeadLinkItem{
			Page:    item.FromPage,
			URL:     item.LinkURL,
			Status:  item.StatusCode,
			Message: item.Err,
		})
	}
	data.DeadLinks = deadLinks
	data.DeadLinkCnt = len(deadLinks)
	return data
}
