package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/public"
	"Blog-Backend/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var collectService = public.NewCollectService()

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

	// 调用geo工具包获取具体信息
	country, region, city := utils.LookupIP(ip)

	// 调用referer解析器获取信息
	referer := c.GetHeader("Referer")
	_, medium, source := utils.ParseReferer(referer)

	info := request.CollectServiceDTO{
		VisitorID: req.VisitorID,
		Path:      req.Path,
		Status:    req.Status,
		Latency:   req.Latency,

		ClientTime: clientTime,
		IP:         ip,
		Country:    country,
		UserAgent:  c.GetHeader("User-Agent"),
		City:       city,
		Region:     region,
		Referer:    referer,
		Medium:     medium,
		Source:     source,
	}

	// 创建上下文
	ctx := c.Request.Context()

	if err := collectService.Collect(ctx, info); err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, gin.H{"status": "ok"})
}
