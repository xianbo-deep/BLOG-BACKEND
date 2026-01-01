package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 除了总日志数，其它都在REDIS拿
func GetDashboardSummary(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := admin.GetDashboardSummary(ctx)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

// 在数据库查前六天，今天的在Redis拿
func GetDashboardTrend(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := admin.GetDashboardTrend(ctx)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)

}

func GetDashboardInsights(c *gin.Context) {
	// 转换类型
	limitstr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitstr)

	res, err := admin.GetDashboardInsights(limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}
