package public

import (
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/dao"
	"Blog-Backend/model"
	"time"
)

func CollectService(info request.CollectServiceDTO) error {

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
	}

	if err := dao.InsertVisitLog(log); err != nil {
		return err
	}

	// 开协程，在redis操作数据
	go func() {
		_ = dao.IncrementPV(info.Path)
		_ = dao.IncrementUV(info.Path, info.VisitorID)
	}()

	return nil

}
