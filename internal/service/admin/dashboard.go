package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
)

// TODO 去REDIS查PV、UV、实时在线人数
func GetDashboardSummary() (response.DashboardSummary, error) {
	var result response.DashboardSummary
	totalLogs, _ := dao.GetTotalLogs()

	result.TotalLogCount = totalLogs

	return result, nil
}

func GetDashboardTrend() ([]response.DashboardTrends, error) {
	// 查过去6天
	history, err := dao.GetHistoryTrends(6)
	if err != nil {
		return nil, err
	}

	// 查今天的
	today, _ := dao.GetTodayPV()

	result := make([]response.DashboardTrends, 0)
	// 倒序查的数据库，需要倒序遍历，让第一天在切片的第一位
	for i := len(history) - 1; i >= 0; i-- {
		result = append(result, history[i])
	}
	// 追加上今天的
	result = append(result, today)

	return result, nil
}

// 查国家分布和访问错误路径的日志
func GetDashboardInsights(limit int) (*response.DashboardInsights, error) {
	geoResult, err := dao.GetGeoDistribution(nil, nil, nil)
	if err != nil {
		return nil, err
	}

	errLogs, err := dao.GetErrorLogs(limit)
	if err != nil {
		return nil, err
	}

	for i := range errLogs {
		t := errLogs[i].Time
		// 转化成时间戳
		errLogs[i].Timestamp = t.Unix()
	}

	return &response.DashboardInsights{
		GeoStats:  geoResult,
		ErrorLogs: errLogs,
	}, nil
}
