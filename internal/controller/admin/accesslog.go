package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAccessLog(c *gin.Context) {
	var req common.PageRequest

	// 查看格式是否正确
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}

	// 调用service
	res, err := admin.GetAccessLog(req)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}

	common.Success(c, res)
}
