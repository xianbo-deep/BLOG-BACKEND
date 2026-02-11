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

type CollectController struct {
	svc *public.CollectService
}

func NewCollectController(svc *public.CollectService) *CollectController {
	return &CollectController{svc: svc}
}

func (ctrl *CollectController) CollectHandler(c *gin.Context) {
	var req request.CollectRequest

	if err := c.ShouldBind(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1001, err.Error())
		return
	}

	// 处理时间
	clientTime := time.UnixMilli(req.Timestamp).UTC()

	// 获取元数据
	meta, _ := GetRequestMeta(c)

	// 限流
	ctx := c.Request.Context()

	ok, err := ctrl.svc.DedupeVisitorPath(ctx, req.VisitorID, req.Path, 2*time.Second)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	// 重复上报
	if !ok {
		common.Success(c, gin.H{"status": "ok"})
		return
	}
	// 调用geo工具包获取具体信息
	geoInfo, _ := utils.LookupIP(meta.IP)
	info := request.CollectServiceDTO{
		VisitorID: req.VisitorID,
		Path:      req.Path,
		Status:    req.Status,
		Latency:   req.Latency,

		ClientTime:  clientTime,
		IP:          meta.IP,
		Country:     geoInfo.CountryZh,
		CountryCode: geoInfo.CountryCode,
		CountryEN:   geoInfo.CountryEn,
		UserAgent:   meta.UserAgent,
		Device:      meta.Device,
		Browser:     meta.Browser,
		OS:          meta.OS,
		City:        geoInfo.CityZh,
		CityEN:      geoInfo.CityEn,
		Region:      geoInfo.RegionZh,
		RegionCode:  geoInfo.RegionCode,
		RegionEN:    geoInfo.RegionEn,
		Referer:     meta.Referer,
		Medium:      meta.Medium,
		Source:      meta.Source,
		Lat:         geoInfo.Lat,
		Lon:         geoInfo.Lon,
	}

	if err := ctrl.svc.Collect(info); err != nil {
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
