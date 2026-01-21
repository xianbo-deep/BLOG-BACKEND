package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var commentService *admin.CommentService

func InitCommentService() {
	commentService = admin.NewCommentService()
}

func GetDiscussionMetric(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	res, err := commentService.GetDiscussionMetric(c, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetDiscussionTrend(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	res, err := commentService.GetDiscussionTrend(c, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetDiscussionNewFeed(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return

	}
	res, err := commentService.GetDiscussionNewFeed(c, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}

func GetDiscussionActiveUser(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	res, err := commentService.GetDiscussionActiveUser(c, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
