package middleware

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/utils"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func HeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var hdr common.RequestHeader
		_ = c.ShouldBindHeader(&hdr)

		ip := GetClientIP(c, hdr)
		device, osName, browser := utils.ParseUA(hdr.UserAgent)
		_, medium, source := utils.ParseReferer(hdr.Referer)

		meta := common.RequestMeta{
			IP:        ip,
			Referer:   hdr.Referer,
			UserAgent: hdr.UserAgent,
			Origin:    hdr.Origin,
			Device:    device,
			OS:        osName,
			Browser:   browser,
			Medium:    medium,
			Source:    source,
		}

		// 存到ctx里面
		c.Set(consts.RequestMetaKey, meta)
		c.Next()
	}
}

func GetClientIP(c *gin.Context, hdr common.RequestHeader) string {
	ip := strings.TrimSpace(hdr.RealIP)
	if ip == "" {
		ip = hdr.GetFirstFowardIP()
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	if ip != "" && net.ParseIP(ip) == nil {
		ip = ""
	}
	return ip
}

func GetRequestMeta(c *gin.Context) (common.RequestMeta, bool) {
	v, ok := c.Get(consts.RequestMetaKey)
	if !ok {
		return common.RequestMeta{}, false
	}
	meta, ok := v.(common.RequestMeta)
	return meta, ok
}
