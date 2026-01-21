package service

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"Blog-Backend/thirdparty/github/query"
	"sort"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
)

/* 返回指标 */

/* 返回最新动态 */
// 组装返回信息
func handleNewFeedRes(allItems []*response.NewFeedItem, limit int) ([]*response.NewFeedItem, error) {
	// 对动态按照时间进行排序
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].Time.After(allItems[j].Time)
	})

	// 截取动态
	if len(allItems) > limit {
		allItems = allItems[:limit]
	}

	for i := range allItems {
		allItems[i].Timestamp = allItems[i].Time.UnixMilli()
	}
	return allItems, nil
}

/* 返回趋势 */
// 添加信息
func addTrendItemReply(trendMap map[string]*response.TrendItem, date time.Time) {
	key := date.Format(consts.DateLayout)
	trendMap[key].TotalReplies++
}

func addTrendItemComment(trendMap map[string]*response.TrendItem, date time.Time) {
	key := date.Format(consts.DateLayout)
	trendMap[key].TotalComments++
}
func addTrendItemReaction(trendMap map[string]*response.TrendItem, date time.Time) {
	key := date.Format(consts.DateLayout)
	trendMap[key].TotalReactions++
}

// 返回前的处理
func handleTrendRes(trendMap map[string]*response.TrendItem, timeRangeDays int) ([]response.TrendItem, error) {
	// 转换成切片
	trends := make([]response.TrendItem, 0, len(trendMap))
	for i := 0; i < timeRangeDays; i++ {
		date := time.Now().AddDate(0, 0, -i)
		key := date.Format(consts.DateLayout)
		trends = append(trends, *trendMap[key])
	}
	return trends, nil
}

/* 返回活跃用户 */
// 更新活跃用户
func addActiveUser(userMap map[string]*response.ActiveUserItem, login, avatar, url string) {
	if login == "" {
		return
	}
	item, ok := userMap[login]
	if !ok {
		item = &response.ActiveUserItem{
			Name:       login,
			TotalFeeds: 0,
			Avatar:     avatar,
			URL:        url,
		}
		userMap[login] = item
	}
	item.TotalFeeds++
}

// 返回前的处理
func handleActiveUserRes(userMap map[string]*response.ActiveUserItem, limit int) ([]response.ActiveUserItem, error) {
	activeusers := make([]response.ActiveUserItem, 0, len(userMap))
	for _, item := range userMap {
		activeusers = append(activeusers, *item)
	}
	// 按活跃度排序
	sort.Slice(activeusers, func(i, j int) bool {
		return activeusers[i].TotalFeeds > activeusers[j].TotalFeeds
	})

	if limit > 0 && limit < len(activeusers) {
		activeusers = activeusers[:limit]
	}
	return activeusers, nil
}

/* 通用模块 */
// 判断是否应该统计
func shouldCount(cufoffTime time.Time, t time.Time) bool {
	return t.After(cufoffTime)
}

// 游标更新
func nextCursor(pi query.PageInfo) *githubv4.String {
	c := pi.EndCursor
	return &c
}

// 合并URL
func concatToUrl(base, url string) string {
	base = strings.TrimSuffix(base, "/")
	url = strings.TrimPrefix(url, "/")
	return base + "/" + url
}
