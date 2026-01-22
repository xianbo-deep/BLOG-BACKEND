package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"
)

type AnalysisService struct {
	dao *dao.AnalysisDao
}

func NewAnalysisService(dao *dao.AnalysisDao) *AnalysisService {
	return &AnalysisService{dao: dao}
}

func (s *AnalysisService) GetAnalysisMetric(days int) (response.AnalysisMetric, error) {
	return s.dao.GetAnalysisMetric(days)
}

func (s *AnalysisService) GetAnalysisTrend(days int) ([]response.AnalysisTrendItem, error) {
	return s.dao.GetAnalysisTrend(days)
}

func (s *AnalysisService) GetAnalysisPathRank(days int) ([]response.AnalysisPathRankItem, error) {
	return s.dao.GetAnalysisPathRank(days)
}

func (s *AnalysisService) GetAnalysisPath(req common.PageRequest, days int) (*common.PageResponse[response.AnalysisPathItem], error) {
	return s.dao.GetAnalysisPath(req, days)
}

func (s *AnalysisService) GetAnalysisPathSource(path string, days int) (response.AnalysisPathItemDetail, error) {
	return s.dao.GetAnalysisPathSource(path, days)
}

func (s *AnalysisService) GetAnalysisPathByQuery(req common.PageRequest, path string, days int) (*common.PageResponse[response.AnalysisPathItem], error) {
	return s.dao.GetAnalysisPathByQuery(req, path, days)
}

func (s *AnalysisService) GetAnalysisPathDetailTrend(path string) ([]response.PathDetailTrendItem, error) {
	return s.dao.GetAnalysisPathDetailTrend(path)
}

func (s *AnalysisService) GetAnalysisPathDetailSource(path string) ([]response.PathDetailSourceItem, error) {
	return s.dao.GetAnalysisPathDetailSource(path)
}
func (s *AnalysisService) GetAnalysisPathDetailDevice(path string) ([]response.PathDetailDeviceItem, error) {
	return s.dao.GetAnalysisPathDetailDevice(path)
}
func (s *AnalysisService) GetAnalysisPathDetailMetric(path string) (response.PathDetailMetric, error) {
	return s.dao.GetAnalysisPathDetailMetric(path)
}
