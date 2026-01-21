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

var analysisService *admin.AnalysisService

func InitAnalysisService(db *gorm.DB) {
	analysisService = admin.NewAnalysisService(db)
}

func GetAnalysisMetrics(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisMetric(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisTrend(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisTrend(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathRank(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathRank(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPath(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPath(common.PageRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}, req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathSource(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathSource(req.Path, req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathByQuery(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := analysisService.GetAnalysisPathByQuery(common.PageRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}, req.Path, req.Days)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
	}
	common.Success(c, res)
}

func GetAnalysisPathDetailTrend(c *gin.Context) {
	path := c.Query("path")
	res, err := analysisService.GetAnalysisPathDetailTrend(path)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
	}
	common.Success(c, res)
}

func GetAnalysisPathDetailMetric(c *gin.Context) {

}

func GetAnalysisPathDetailSource(c *gin.Context) {

}

func GetAnalysisPathDetailDevice(c *gin.Context) {

}
