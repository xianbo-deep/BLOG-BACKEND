package public

import (
	"Blog-Backend/core"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CollectHandler(c *gin.Context) {
	/* 解析参数 */
	path := c.Query("path")
	country := c.GetHeader("x-vercel-ip-country")
	city := c.GetHeader("x-vercel-ip-city")
	ua := c.Request.UserAgent()
	ip := c.GetHeader("x-real-ip")
	region := c.GetHeader("x-vercel-ip-country-region")
	status := c.GetHeader("status")
	if ip == "" {
		ip = c.GetHeader("x-forwarded-for")
	}
	referer := c.GetHeader("referer")

	/* 解析时间 */
	ctime := c.Query("client-time")
	var clientTime time.Time
	if ctime != "" {
		clientTime, _ = time.Parse(time.RFC3339, ctime) // 前端传 ISO 格式
	}

	/* 插入数据 */
	if core.RDB != nil {

	}

	if core.DB != nil {
		log := core.VisitLog{
			VisitTime:  time.Now(),
			ClientTime: clientTime,
			Path:       path,
			Country:    country,
			City:       city,
			UserAgent:  ua,
			IP:         ip,
			Region:     region,
			Referer:    referer,
			Status:     status,
		}

		core.DB.Create(&log)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
