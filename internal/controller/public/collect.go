package public

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/public"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CollectHandler(c *gin.Context) {
	var req request.CollectRequest

	if err := c.ShouldBind(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1001, err.Error())
		return
	}

	// 处理时间
	clientTime := time.UnixMilli(req.Timestamp)

	// 处理ip
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
		if ip != "" {
			ip = strings.Split(ip, ",")[0]
			ip = strings.TrimSpace(ip)
		}
	}
	if ip == "" {
		ip = c.ClientIP()
	}

	info := request.CollectServiceDTO{
		VisitorID: req.VisitorID,
		Path:      req.Path,
		Status:    req.Status,
		Latency:   req.Latency,

		ClientTime: clientTime,
		IP:         ip,
		Country:    c.GetHeader("x-vercel-ip-country"),
		UserAgent:  c.GetHeader("User-Agent"),
		City:       c.GetHeader("x-vercel-ip-city"),
		Region:     c.GetHeader("x-vercel-ip-country-region"),
		Referer:    c.GetHeader("Referer"),
	}

	if err := public.CollectService(info); err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, gin.H{"status": "ok"})
}
