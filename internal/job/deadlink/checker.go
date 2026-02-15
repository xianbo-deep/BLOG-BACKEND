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
	workDir, cleanup, err := c.prepareWorkTree()
	if err != nil {
		return Summary{}, nil, err
	}
	defer cleanup()

	// 获取存储md文件的目录
	docsPath := filepath.Join(workDir, c.cfg.DocsDir)

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

// 准备worktree
func (c *Checker) prepareWorkTree() (string, func(), error) {
	// 获取必须的参数
	cache := c.cfg.CacheRepoDir
	if cache == "" {
		cache = CacheRepoDir
	}
	branch := c.cfg.Branch
	if branch == "" {
		branch = Branch
	}

	// 缓存仓库不存在
	if _, err := os.Stat(cache); err != nil {
		if err := os.MkdirAll(filepath.Dir(cache), 0755); err != nil {
			return "", nil, err
		}
		if err := c.gitCloneMirror(cache, branch); err != nil {
			return "", nil, err
		}
	}

	// fetch获取最新的变更文件
	if err := c.gitFetch(cache, branch); err != nil {
		return "", nil, err
	}

	wt, err := os.MkdirTemp("", "deadlink-wt-*")
	if err != nil {
		return "", nil, err
	}

	// 创建分离工作树
	if err := c.gitCmd(2*time.Minute, cache,
		"worktree", "add", "--detach", wt, "origin/"+branch,
	); err != nil {
		_ = os.RemoveAll(wt)
		return "", nil, err
	}

	cleanup := func() {
		// Git层面移除工作树
		_ = c.gitCmd(2*time.Minute, cache, "worktree", "remove", "--force", wt)
		// 物理删除临时目录
		_ = os.RemoveAll(wt)
	}
	return wt, cleanup, nil
}

// 克隆仓库
func (c *Checker) gitCloneMirror(cache, branch string) error {
	return c.gitCmd(5*time.Minute, "", // 注意：clone 不用 -C
		"-c", "http.version=HTTP/1.1",
		"clone",
		"--mirror",
		"--single-branch",
		"--branch", branch,
		c.cfg.RepoURL,
		cache,
	)
}

// fetch获取最新变更文件
func (c *Checker) gitFetch(cache, branch string) error {
	return c.gitCmd(3*time.Minute, cache,
		"-c", "http.version=HTTP/1.1",
		"fetch",
		"--prune",
		"origin",
		"+refs/heads/"+branch+":refs/remotes/origin/"+branch)
}

// 命令行函数
func (c *Checker) gitCmd(timeout time.Duration, dir string, args ...string) error {
	ctx, cancel := consts.GetTimeoutContext(context.Background(), timeout)
	defer cancel()

	finalArgs := args

	if dir != "" {
		finalArgs = append([]string{"-C", dir}, args...)
	}

	// 注入参数
	cmd := exec.CommandContext(ctx, "git", finalArgs...)
	// 获取环境变量
	cmd.Env = c.gitEnv()
	// 获取输出
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("git操作超时 :%s", string(out))
	}

	if err != nil {
		return fmt.Errorf("git操作失败 :%s", string(out))
	}
	return nil
}

// 继承环境变量并注入代理配置
func (c *Checker) gitEnv() []string {
	env := os.Environ()

	proxy := strings.TrimSpace(c.cfg.ProxyHTTP)
	if proxy == "" {
		proxy = ProxyHTTP
	}

	env = append(env, "HTTP_PROXY="+proxy, "HTTPS_PROXY="+proxy)

	return env
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
	reURL := regexp.MustCompile(`https?://[A-Za-z0-9][A-Za-z0-9.-]*(?:\.[A-Za-z0-9.-]+)+(?:/[^\s\)\]\}>"']*)?`)
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
