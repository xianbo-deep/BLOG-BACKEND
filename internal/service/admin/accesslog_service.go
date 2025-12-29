package admin

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"Blog-Backend/utils"
)

func GetAccessLog(req common.PageRequest) (*common.PageResponse[response.AccessLog], error) {
	db := core.DB.Order("visit_time desc")

	// 查的时候用了实体类
	pageResult, err := utils.Paginate[model.VisitLog](db, req)

	if err != nil {
		return nil, err
	}

	// 进行转换，提取有用的信息
	var dtoList []response.AccessLog

	for _, v := range pageResult.List {
		dtoList = append(dtoList, response.AccessLog{
			Path:       v.Path,
			VisitTime:  v.VisitTime,
			IP:         v.IP,
			ClientTime: v.ClientTime,
			UserAgent:  v.UserAgent,
			Referer:    v.Referer,
			Country:    v.Country,
			City:       v.City,
			Region:     v.Region,
			Status:     v.Status,
		})
	}

	return &common.PageResponse[response.AccessLog]{
		List:      dtoList,
		Total:     pageResult.Total,
		Page:      pageResult.Page,
		PageSize:  pageResult.PageSize,
		TotalPage: pageResult.TotalPage,
	}, nil
}
