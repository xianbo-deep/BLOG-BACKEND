package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	svc *admin.DashboardService
}

func NewDashboardController(svc *admin.DashboardService) *DashboardController {
	return &DashboardController{svc: svc}
}

// 除了总日志数，其它都在REDIS拿
func (ctrl *DashboardController) GetDashboardSummary(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := ctrl.svc.GetDashboardSummary(ctx)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

// 在数据库查前六天，今天的在Redis拿
func (ctrl *DashboardController) GetDashboardTrend(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := ctrl.svc.GetDashboardTrend(ctx)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)

}

func (ctrl *DashboardController) GetDashboardInsights(c *gin.Context) {
	// 转换类型
	limitstr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitstr)

	res, err := ctrl.svc.GetDashboardInsights(limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
