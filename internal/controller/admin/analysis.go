package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

var analysisService = admin.NewAnalysisService()

func GetAnalysisMetrics(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisMetric(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisTrend(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisTrend(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathRank(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathRank(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPath(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPath(common.PageRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}, req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathDetail(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathDetail(req.Path, req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathByQuery(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathByQuery(common.PageRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}, req.Path, req.Days)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
	}
	common.Success(c, res)
}
