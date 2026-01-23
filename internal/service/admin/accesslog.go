package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"Blog-Backend/utils"
	"os"
	"strconv"

	"gorm.io/gorm"
)

type AccessLogService struct {
	db *gorm.DB
}

func NewAccessLogService(db *gorm.DB) *AccessLogService {
	return &AccessLogService{db: db}
}

func (s *AccessLogService) GetAccessLogByQuery(req request.AccessLogRequest) (*common.PageResponse[response.AccessLog], error) {
	db := s.db.Model(&model.VisitLog{}).Order("visit_time desc")

	if req.KeyWord != "" {
		keyword := "%" + req.KeyWord + "%"
		db = db.Where("path LIKE ? or ip LIKE ? or visitor_id LIKE ?", keyword, keyword, keyword)
	}

	if req.Status != "" {
		if statusInt, err := strconv.Atoi(req.Status); err == nil {
			db = db.Where("status = ?", statusInt)
		}
	}

	if req.Latency != 0 {
		db = db.Where("latency > ?", req.Latency)
	}

	// 分页查
	pageResult, err := utils.Paginate[model.VisitLog](db, req.PageRequest)
	if err != nil {
		return nil, err
	}
	var dtoList []response.AccessLog
	for _, v := range pageResult.List {
		dtoList = append(dtoList, response.AccessLog{
			VisitorID:       v.VisitorID,
			Path:            os.Getenv(consts.EnvBaseURL) + v.Path,
			VisitTime:       consts.TransferTimeByLoc(v.VisitTime).Format(consts.TimeLayout),
			IP:              v.IP,
			ClientTime:      consts.TransferTimeByLoc(v.ClientTime).Format(consts.TimeLayout),
			UserAgent:       v.UserAgent,
			Referer:         v.Referer,
			Country:         v.Country,
			City:            v.City,
			Region:          v.Region,
			Status:          v.Status,
			Latency:         v.Latency,
			Browser:         v.Browser,
			Device:          v.Device,
			OS:              v.OS,
			Source:          v.Source,
			Medium:          v.Medium,
			VistitTimestamp: consts.TransferTimeToTimestamp(v.VisitTime),
			ClientTimestamp: consts.TransferTimeToTimestamp(v.ClientTime),
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
