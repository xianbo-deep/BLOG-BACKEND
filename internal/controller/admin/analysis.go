package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

var analysisService = admin.NewAnalysisService()

func GetAnalysisMetrics(c *gin.Context) {
	res, err := analysisService.GetAnalysisMetric()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisTrend(c *gin.Context) {
	res, err := analysisService.GetAnalysisTrend()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathRank(c *gin.Context) {
	res, err := analysisService.GetAnalysisPathRank()
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

	res, err := analysisService.GetAnalysisPath(req)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetAnalysisPathDetail(c *gin.Context) {

}

func GetAnalysisPathByQuery(c *gin.Context) {

}
