package dao

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/response"
	"Blog-Backend/model"
	"time"

	"gorm.io/gorm"
)

type AnalysisDao struct {
	db *gorm.DB
}

func NewAnalysisDao(db *gorm.DB) *AnalysisDao {
	return &AnalysisDao{db: db}
}

func (d *AnalysisDao) GetAnalysisMetric(days int) (response.AnalysisMetric, error) {
	cutoffTime := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay).AddDate(0, 0, -days)
	var res response.AnalysisMetric
	err := d.db.Model(&model.VisitLog{}).
		Select(`
				count(*) as total_pv,
				count(distinct visitor_id) as total_uv,
				avg(latency)::bigint as avg_latency`). // TODO 看看这个用法
		Where("visit_time > ?", cutoffTime).
		Scan(&res).
		Error
	if err != nil {
		return response.AnalysisMetric{}, err
	}

	var hotpage response.HotPageResult

	err = d.db.Model(&model.VisitLog{}).
		Select("path,count(*) as pv").
		Group("path").
		Order("pv desc").
		Limit(1).
		Where("visit_time > ?", cutoffTime).
		Scan(&hotpage).
		Error
	if err != nil {
		return response.AnalysisMetric{}, err
	}

	res.HotPage = hotpage.Path
	res.HotPagePV = hotpage.PV

	return res, nil
}

func (d *AnalysisDao) GetAnalysisTrend(days int) ([]response.AnalysisTrendItem, error) {
	cutoffTime := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay).AddDate(0, 0, -days)
	var res []response.AnalysisTrendItem
	err := d.db.Model(&model.VisitLog{}).
		Select("date(visit_time AT TIME ZONE 'Asia/Shanghai') as date,count(*) as pv,count(distinct visitor_id) as uv").
		Where("visit_time > ?", cutoffTime).
		Group("date(visit_time AT TIME ZONE 'Asia/Shanghai')").
		Order("date asc").
		Scan(&res).
		Error
	if err != nil {
		return nil, err
	}
	for i := range res {
		res[i].Timestamp = consts.TransferTimeToTimestamp(res[i].Date)
	}
	return res, nil
}

func (d *AnalysisDao) GetAnalysisPathRank(days int) ([]response.AnalysisPathRankItem, error) {
	cutoffTime := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay).AddDate(0, 0, -days)
	var res []response.AnalysisPathRankItem
	err := d.db.Model(&model.VisitLog{}).
		Select("path,count(*) as pv").
		Where("visit_time > ?", cutoffTime).
		Order("pv desc").
		Group("path").
		Limit(10).
		Scan(&res).
		Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *AnalysisDao) GetAnalysisPath(req common.PageRequest, days int) (*common.PageResponse[response.AnalysisPathItem], error) {
	cutoffTime := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay).AddDate(0, 0, -days)
	db := d.db.Where("visit_time >= ?", cutoffTime)

	// 这里不能用分页插件，因为返回的不是原始数据模型的数据，是其聚合之后的数据
	var total int64
	page := req.GetPage()
	pageSize := req.GetPageSize()
	offset := (page - 1) * pageSize

	// 查总数
	if err := db.Model(&model.VisitLog{}).
		Distinct("path").
		Count(&total).Error; err != nil {
		return nil, err
	}

	// 查最后的结果
	var res []response.AnalysisPathItem
	err := db.Model(&model.VisitLog{}).
		Select("path,count(*) as pv,count(distinct visitor_id) as uv,avg(latency)::bigint as avg_latency").
		Group("path").
		Order("pv desc").
		Offset(offset).
		Limit(pageSize).
		Scan(&res).Error
	if err != nil {
		return nil, err
	}

	totalPage := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &common.PageResponse[response.AnalysisPathItem]{
		List:      res,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (d *AnalysisDao) GetAnalysisPathByQuery(req common.PageRequest, path string, days int) (*common.PageResponse[response.AnalysisPathItem], error) {
	cutoffTime := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay).AddDate(0, 0, -days)
	db := d.db.Where("visit_time >= ? and path like ?", cutoffTime, "%"+path+"%")

	// 这里不能用分页插件，因为返回的不是原始数据模型的数据，是其聚合之后的数据
	var total int64
	page := req.GetPage()
	pageSize := req.GetPageSize()
	offset := (page - 1) * pageSize

	// 查总数
	if err := db.Model(&model.VisitLog{}).
		Distinct("path").
		Count(&total).Error; err != nil {
		return nil, err
	}

	// 查最后的结果
	var res []response.AnalysisPathItem
	err := db.Model(&model.VisitLog{}).
		Select("path,count(*) as pv,count(distinct visitor_id) as uv,avg(latency)::bigint as avg_latency").
		Group("path").
		Order("pv desc").
		Offset(offset).
		Limit(pageSize).
		Scan(&res).Error
	if err != nil {
		return nil, err
	}

	totalPage := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &common.PageResponse[response.AnalysisPathItem]{
		List:      res,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (d *AnalysisDao) GetAnalysisPathSource(path string, days int) (response.AnalysisPathItemDetail, error) {
	cutoffTime := consts.TransferTimeByLoc(time.Now()).Truncate(consts.TimeRangeDay).AddDate(0, 0, -days)
	db := d.db.Model(&model.VisitLog{}).Where("visit_time >= ? and path like ?", cutoffTime, "%"+path+"%")
	var res response.AnalysisPathItemDetail

	var totalPV int64
	if err := db.Count(&totalPV).Error; err != nil {
		return res, err
	}

	// 没有访问记录，直接返回空
	if totalPV == 0 {
		return res, nil
	}

	// 统计country
	var countries []struct {
		Country string
		Count   int64
	}

	err := db.Session(&gorm.Session{}).
		Select("country,count(*) as count").
		Group("country").
		Order("count desc").
		Limit(3).
		Scan(&countries).Error
	if err != nil {
		return res, err
	}

	// 填充数据
	for _, r := range countries {
		percent := int64(float64(r.Count) / float64(totalPV) * 100)
		res.Countries = append(res.Countries, response.AnalysisPathItemCountry{
			Country: r.Country,
			Percent: percent,
		})
	}

	// 统计device
	var devices []struct {
		Device string
		Count  int64
	}

	err = db.Session(&gorm.Session{}).
		Select("device, count(*) as count").
		Group("device").
		Order("count desc").
		Limit(3).
		Scan(&devices).Error
	if err != nil {
		return res, err
	}

	// 填充数据
	for _, c := range devices {
		percent := int64(float64(c.Count) / float64(totalPV) * 100)
		res.Devices = append(res.Devices, response.AnalysisPathItemDevice{
			Device:  c.Device,
			Percent: percent,
		})
	}
	// 返回数据
	res.Path = path
	return res, nil
}

func (d *AnalysisDao) GetAnalysisPathDetailTrend(path string) ([]response.PathDetailTrendItem, error) {
	var res []response.PathDetailTrendItem
	db := d.db.Model(&model.VisitLog{})

	startTime := consts.TransferTimeByLoc(time.Now()).Add(-consts.TimeRangeDay)

	err := db.Select("date_trunc('hour',visit_time AT TIME ZONE 'Asia/Shanghai') as date,count (*) as pv,count(distinct visitor_id) as uv").
		Where("visit_time > ? and path = ?", startTime, path).
		Group("date_trunc('hour',visit_time AT TIME ZONE 'Asia/Shanghai')").
		Order("date asc").
		Scan(&res).Error

	if err != nil {
		return nil, err
	}
	for i := range res {
		res[i].Date = consts.TransferTimeByLoc(res[i].Date)
		res[i].Timestamp = consts.TransferTimeToTimestamp(res[i].Date)
	}
	return res, nil
}

func (d *AnalysisDao) GetAnalysisPathDetailSource(path string) ([]response.PathDetailCountryItem, error) {
	var res []response.PathDetailCountryItem
	db := d.db.Model(&model.VisitLog{})

	err := db.Select("country,count(*) as count").
		Where("path = ?", path).
		Group("country").
		Order("count desc").
		Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *AnalysisDao) GetAnalysisPathDetailDevice(path string) ([]response.PathDetailDeviceItem, error) {
	var res []response.PathDetailDeviceItem
	db := d.db.Model(&model.VisitLog{})

	err := db.Select("device,coalesce(count(*),0) as count").
		Where("path = ?", path).
		Group("device").
		Order("count desc").
		Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (d *AnalysisDao) GetAnalysisPathDetailMetric(path string) (response.PathDetailMetric, error) {
	db := d.db.Model(&model.VisitLog{})
	var res response.PathDetailMetric
	err := db.Select("count(distinct visitor_id) as uv,count(*) as pv").
		Where("path = ?", path).
		Scan(&res).Error
	if err != nil {
		return res, err
	}
	return res, nil
}
