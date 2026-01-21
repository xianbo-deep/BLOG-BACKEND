package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var accesslogService *admin.AccessLogService

func InitAccessLogService(db *gorm.DB) {
	accesslogService = admin.NewAccessLogService(db)
}

func GetAccessLog(c *gin.Context) {
	var req common.PageRequest

	// 查看格式是否正确
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}

	// 调用service
	res, err := accesslogService.GetAccessLog(req)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}

	common.Success(c, res)
}

func GetAccessLogByQuery(c *gin.Context) {
	var req request.AccessLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}

	// 调用service
	res, err := accesslogService.GetAccessLogByQuery(req)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
