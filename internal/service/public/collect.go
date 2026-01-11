package public

import (
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/dao"
	"Blog-Backend/model"
	"context"
	"time"
)

type CollectService struct{}

func NewCollectService() *CollectService {
	return &CollectService{}
}

func (s *CollectService) Collect(ctx context.Context, info request.CollectServiceDTO) error {

	log := model.VisitLog{
		VisitTime:  time.Now(),
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
		_ = dao.IncrementPV(ctx, info.Path)
		_ = dao.IncrementUV(ctx, info.Path, info.VisitorID)
		_ = dao.RecordOnline(ctx, info.VisitorID)
		_ = dao.RecordLatency(ctx, info.Path, info.Latency)
	}()

	return nil

}
