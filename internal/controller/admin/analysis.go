package admin

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

var analysisService = admin.NewAnalysisService(core.DB)

func GetTotalPagesData(c *gin.Context) {
	var req common.PageRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}

	res, err := analysisService.GetTotalPagesData(req)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
	}

	common.Success(c, res)

}

func GetTodayPagesData(c *gin.Context) {
	var req common.PageRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}

	res, err := admin.GetTodayPagesData(req)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
	}

	common.Success(c, res)
}
