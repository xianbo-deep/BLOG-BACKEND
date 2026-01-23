package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnalysisController struct {
	svc *admin.AnalysisService
}

func NewAnalysisController(svc *admin.AnalysisService) *AnalysisController {
	return &AnalysisController{svc: svc}
}

func (ctrl *AnalysisController) GetAnalysisMetrics(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := ctrl.svc.GetAnalysisMetric(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisTrend(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := ctrl.svc.GetAnalysisTrend(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathRank(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := ctrl.svc.GetAnalysisPathRank(req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPath(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := ctrl.svc.GetAnalysisPath(common.PageRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}, req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathSource(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := ctrl.svc.GetAnalysisPathSource(req.Path, req.Days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathByQuery(c *gin.Context) {
	var req request.AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}
	res, err := ctrl.svc.GetAnalysisPathByQuery(common.PageRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}, req.Path, req.Days)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathDetailTrend(c *gin.Context) {
	path := c.Query("path")
	res, err := ctrl.svc.GetAnalysisPathDetailTrend(path)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathDetailMetric(c *gin.Context) {
	path := c.Query("path")
	res, err := ctrl.svc.GetAnalysisPathDetailMetric(path)

	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}

	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathDetailSource(c *gin.Context) {
	path := c.Query("path")
	res, err := ctrl.svc.GetAnalysisPathDetailSource(path)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *AnalysisController) GetAnalysisPathDetailDevice(c *gin.Context) {
	path := c.Query("path")
	res, err := ctrl.svc.GetAnalysisPathDetailDevice(path)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
