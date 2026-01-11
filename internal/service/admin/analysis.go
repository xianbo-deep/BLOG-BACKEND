package admin

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"

	"gorm.io/gorm"
)

type AnalysisService struct {
	db *gorm.DB
}

func NewAnalysisService() *AnalysisService {
	return &AnalysisService{db: core.DB}
}

func (s *AnalysisService) GetAnalysisMetric(days int) (response.AnalysisMetric, error) {
	return dao.GetAnalysisMetric(days)
}

func (s *AnalysisService) GetAnalysisTrend(days int) ([]response.AnalysisTrendItem, error) {
	return dao.GetAnalysisTrend(days)
}

func (s *AnalysisService) GetAnalysisPathRank(days int) ([]response.AnalysisPathRankItem, error) {
	return dao.GetAnalysisPathRank(days)
}

func (s *AnalysisService) GetAnalysisPath(req common.PageRequest, days int) (*common.PageResponse[response.AnalysisPathItem], error) {
	return dao.GetAnalysisPath(req, days)
}

func (s *AnalysisService) GetAnalysisPathSource(path string, days int) (response.AnalysisPathItemDetail, error) {
	return dao.GetAnalysisPathSource(path, days)
}

func (s *AnalysisService) GetAnalysisPathByQuery(req common.PageRequest, path string, days int) (*common.PageResponse[response.AnalysisPathItem], error) {
	return dao.GetAnalysisPathByQuery(req, path, days)
}

func (s *AnalysisService) GetAnalysisPathDetailTrend(path string) ([]response.PathDetailTrendItem, error) {
	return dao.GetAnalysisPathDetailTrend(path)
}

func (s *AnalysisService) GetAnalysisPathDetailSource(path string) ([]response.PathDetailSourceItem, error) {
	return dao.GetAnalysisPathDetailSource(path)
}
func (s *AnalysisService) GetAnalysisPathDetailDevice(path string) ([]response.PathDetailDeviceItem, error) {
	return dao.GetAnalysisPathDetailDevice(path)
}
func (s *AnalysisService) GetAnalysisPathDetailMetric(path string) (response.PathDetailMetric, error) {
	return dao.GetAnalysisPathDetailMetric(path)
}
