package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
	"context"
)

type DashboardService struct {
	dao *dao.DashboardDao
}

func NewDashboardService(dao *dao.DashboardDao) *DashboardService {
	return &DashboardService{dao: dao}
}

// 去REDIS查PV、UV、实时在线人数
func (s *DashboardService) GetDashboardSummary(ctx context.Context) (response.DashboardSummary, error) {
	var result response.DashboardSummary
	// 获取总日志数
	totalLogs, _ := s.dao.GetTotalLogs()
	result.TotalLogCount = totalLogs
	// 获取在线人数
	count, _ := s.dao.GetOnlineCount(ctx)
	result.OnlineCount = count
	// 获取PV和UV
	UV, PV, _ := s.dao.GetTodayPVUV(ctx)
	result.UV = UV
	result.PV = PV
	// 获取前一天的PV和UV
	res, err := s.dao.GetLastDayPVUV()
	if err != nil {
		return result, err
	}
	if res.PV > 0 {
		result.PvPercent = (float64(PV) - float64(res.PV)) / float64(res.PV)
	} else {
		if PV > 0 {
			result.PvPercent = 1.0
		} else {
			result.PvPercent = 0.0
		}
	}
	if res.UV > 0 {
		result.UVPercent = (float64(UV) - float64(res.UV)) / float64(res.UV)
	} else {
		if UV > 0 {
			result.UVPercent = 1.0
		} else {
			result.UVPercent = 0.0
		}
	}
	return result, nil
}

// 查博客总趋势
func (s *DashboardService) GetDashboardTrend(ctx context.Context) ([]response.DashboardTrends, error) {
	// 查过去6天
	history, err := s.dao.GetHistoryTrends(6)
	if err != nil {
		return nil, err
	}

	// 查今天的
	today, _ := s.dao.GetTodayPV(ctx)

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
func (s *DashboardService) GetDashboardInsights(limit int) (*response.DashboardInsights, error) {
	geoResult, err := s.dao.GetGeoDistribution(nil, nil, nil)
	if err != nil {
		return nil, err
	}

	errLogs, err := s.dao.GetErrorLogs(limit)
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
