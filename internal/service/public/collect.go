package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/dao"
	"Blog-Backend/model"
	"context"
)

type CollectService struct{}

func NewCollectService() *CollectService {
	return &CollectService{}
}

func (s *CollectService) Collect(ctx context.Context, info request.CollectServiceDTO) error {

	log := model.VisitLog{
		VisitTime:  consts.GetCurrentUTCTime(),
		ClientTime: info.ClientTime,
		Path:       info.Path,
		Country:    info.Country,
		City:       info.City,
		UserAgent:  info.UserAgent,
		IP:         info.IP,
		Region:     info.Region,
		Referer:    info.Referer,
		Status:     info.Status,
		VisitorID:  info.VisitorID,
		Latency:    info.Latency,
		Medium:     info.Medium,
		Source:     info.Source,
		Device:     info.Device,
		OS:         info.OS,
		Browser:    info.Browser,
	}

	if err := dao.InsertVisitLog(log); err != nil {
		return err
	}

	// 开协程，在redis操作数据
	go func() {
		bg, cancel := consts.GetTimeoutContext(context.Background(), consts.RedisOperationTimeout)
		defer cancel()
		_ = dao.IncrementPV(bg, info.Path)
		_ = dao.IncrementUV(bg, info.Path, info.VisitorID)
		_ = dao.RecordOnline(bg, info.VisitorID)
		_ = dao.RecordLatency(bg, info.Path, info.Latency)
	}()

	return nil

}
