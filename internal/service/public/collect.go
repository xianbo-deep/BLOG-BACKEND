package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/dao"
	"Blog-Backend/internal/ws"
	"Blog-Backend/model"
	"context"
)

type CollectService struct {
	dao *dao.CollectDao
	hub *ws.Hub
}

func NewCollectService(dao *dao.CollectDao, hub *ws.Hub) *CollectService {
	return &CollectService{dao: dao, hub: hub}
}

func (s *CollectService) Collect(info request.CollectServiceDTO) error {

	log := model.VisitLog{
		VisitTime:   consts.GetCurrentUTCTime(),
		ClientTime:  info.ClientTime,
		Path:        info.Path,
		Country:     info.Country,
		CountryCode: info.CountryCode,
		CountryEN:   info.CountryEN,
		CityEN:      info.CityEN,
		RegionCode:  info.RegionCode,
		RegionEN:    info.RegionEN,
		City:        info.City,
		UserAgent:   info.UserAgent,
		IP:          info.IP,
		Region:      info.Region,
		Referer:     info.Referer,
		Status:      info.Status,
		VisitorID:   info.VisitorID,
		Latency:     info.Latency,
		Medium:      info.Medium,
		Source:      info.Source,
		Device:      info.Device,
		OS:          info.OS,
		Browser:     info.Browser,
		Lat:         info.Lat,
		Lon:         info.Lon,
	}

	if err := s.dao.InsertVisitLog(log); err != nil {
		return err
	}

	// 开协程，在redis操作数据
	go func() {
		bg, cancel := consts.GetTimeoutContext(context.Background(), consts.RedisOperationTimeout)
		defer cancel()
		_ = s.dao.IncrementPV(bg, info.Path)
		_ = s.dao.IncrementUV(bg, info.Path, info.VisitorID)
		_ = s.dao.RecordOnline(bg, info.VisitorID)
		_ = s.dao.RecordLatency(bg, info.Path, info.Latency)

		// 广播
		if s.hub != nil {

		}
	}()

	return nil

}
