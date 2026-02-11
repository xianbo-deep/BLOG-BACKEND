package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PerformanceController struct {
	svc *admin.PerformanceService
}

func NewPerformanceController(svc *admin.PerformanceService) *PerformanceController {
	return &PerformanceController{svc: svc}
}

func (ctrl *PerformanceController) GetAverageDelay(c *gin.Context) {
	res, err := ctrl.svc.GetAverageDelay()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *PerformanceController) GetSlowPages(c *gin.Context) {
	ctx := c.Request.Context()
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	res, err := ctrl.svc.GetSlowPages(ctx, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
