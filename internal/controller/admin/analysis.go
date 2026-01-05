package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var analysisService = admin.NewAnalysisService()

// TODO 后面记得重构POST请求的代码 只用GET即可
func GetAnalysisMetrics(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 1000, e.Error())
		return
	}
	res, err := analysisService.GetAnalysisMetric(days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisTrend(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 1000, e.Error())
		return
	}
	res, err := analysisService.GetAnalysisTrend(days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathRank(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 1000, e.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathRank(days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPath(c *gin.Context) {
	var req common.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 1000, e.Error())
		return
	}
	res, err := analysisService.GetAnalysisPath(req, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathDetail(c *gin.Context) {
	daysStr := c.Query("days")
	path := c.Query("path")

	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 1000, e.Error())
	}
	res, err := analysisService.GetAnalysisPathDetail(path, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
	}
	common.Success(c, res)
}

func GetAnalysisPathByQuery(c *gin.Context) {
	var req common.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}

	daysStr := c.Query("days")
	path := c.Query("path")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 1001, e.Error())
	}
	res, err := analysisService.GetAnalysisPathByQuery(req, path, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
	}
	common.Success(c, res)
}
