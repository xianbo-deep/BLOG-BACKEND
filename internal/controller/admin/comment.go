package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	svc *admin.CommentService
}

func NewCommentController(svc *admin.CommentService) *CommentController {
	return &CommentController{svc: svc}
}

func (ctrl *CommentController) GetDiscussionMetric(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	res, err := ctrl.svc.GetDiscussionMetric(c, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *CommentController) GetDiscussionTrend(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	res, err := ctrl.svc.GetDiscussionTrend(c, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *CommentController) GetDiscussionNewFeed(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return

	}
	res, err := ctrl.svc.GetDiscussionNewFeed(c, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func (ctrl *CommentController) GetDiscussionActiveUser(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	res, err := ctrl.svc.GetDiscussionActiveUser(c, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
