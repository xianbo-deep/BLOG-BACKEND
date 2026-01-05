package admin

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/dao"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalysisService struct {
	db *gorm.DB
}

func NewAnalysisService() *AnalysisService {
	return &AnalysisService{db: core.DB}
}

func (s *AnalysisService) GetAnalysisMetric() (response.Metric, error) {
	return dao.GetAnalysisMetric()
}

func (s *AnalysisService) GetAnalysisTrend() ([]response.AnalysisTrendItem, error) {
	return dao.GetAnalysisTrend()
}

func (s *AnalysisService) GetAnalysisPathRank() ([]response.AnalysisPathRankItem, error) {

}

func (s *AnalysisService) GetAnalysisPath(req common.PageRequest) ([]common.PageResponse[response.AnalysisPathItem], error) {

}

func (s *AnalysisService) GetAnalysisPathDetail(path string) (response.AnalysisPathItemDetail, error) {
}

func (s *AnalysisService) GetAnalysisPathByQuery(req common.PageRequest) ([]common.PageResponse[response.AnalysisPathItem], error) {
}
