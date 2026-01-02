package admin

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"Blog-Backend/model"
	"Blog-Backend/utils"

	"gorm.io/gorm"
)

type AnalysisService struct {
	db *gorm.DB
}

func NewAnalysisService(db *gorm.DB) *AnalysisService {
	return &AnalysisService{db: db}
}

func (s *AnalysisService) GetTotalPagesData(req common.PageRequest) (*common.PageResponse[model.PathStatRes], error) {
	db := s.db

	pageResult, err := utils.Paginate[model.PathStatRes](db, req)

	if err != nil {
		return nil, err
	}

	var dtoList []model.PathStatRes

	for _, v := range pageResult.List {
		dtoList = append(dtoList, model.PathStatRes{
			Path:       v.Path,
			PV:         v.PV,
			UV:         v.UV,
			AvgLatency: v.AvgLatency,
		})
	}

	return &common.PageResponse[model.PathStatRes]{
		List:      dtoList,
		Total:     pageResult.Total,
		Page:      pageResult.Page,
		PageSize:  pageResult.PageSize,
		TotalPage: pageResult.TotalPage,
	}, nil

}

func GetTodayPagesData(req common.PageRequest) (*common.PageResponse[model.PathStatRes], error) {

}
