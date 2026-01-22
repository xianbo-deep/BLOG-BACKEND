package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccessLogController struct {
	svc *admin.AccessLogService
}

func NewAccessLogController(svc *admin.AccessLogService) *AccessLogController {
	return &AccessLogController{svc: svc}
}

func (ctrl *AccessLogController) GetAccessLogByQuery(c *gin.Context) {
	var req request.AccessLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}

	// 调用service
	res, err := ctrl.svc.GetAccessLogByQuery(req)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
