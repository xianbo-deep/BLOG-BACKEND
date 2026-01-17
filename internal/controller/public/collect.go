package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/public"
	"Blog-Backend/utils"
	"net/http"
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
	clientTime := time.UnixMilli(req.Timestamp).UTC()

	// 获取元数据
	meta, _ := GetRequestMeta(c)

	// 调用geo工具包获取具体信息
	country, region, city := utils.LookupIP(meta.IP)

	info := request.CollectServiceDTO{
		VisitorID: req.VisitorID,
		Path:      req.Path,
		Status:    req.Status,
		Latency:   req.Latency,

		ClientTime: clientTime,
		IP:         meta.IP,
		Country:    country,
		UserAgent:  meta.UserAgent,
		Device:     meta.Device,
		Browser:    meta.Browser,
		OS:         meta.OS,
		City:       city,
		Region:     region,
		Referer:    meta.Referer,
		Medium:     meta.Medium,
		Source:     meta.Source,
	}

	// 创建上下文
	ctx := c.Request.Context()

	if err := collectService.Collect(ctx, info); err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, gin.H{"status": "ok"})
}

func GetRequestMeta(c *gin.Context) (common.RequestMeta, bool) {
	v, ok := c.Get(consts.RequestMetaKey)
	if !ok {
		return common.RequestMeta{}, false
	}
	meta, ok := v.(common.RequestMeta)
	return meta, ok
}
