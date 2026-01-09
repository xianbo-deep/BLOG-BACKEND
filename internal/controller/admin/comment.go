package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/thirdparty/github"
	"Blog-Backend/thirdparty/github/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var discussionService = service.NewDiscussionService(github.NewClient())

func GetDiscussionMetric(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 2000, e.Error())
		return
	}
	res, err := discussionService.GetTotalMetric(c, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 1000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetDiscussionTrend(c *gin.Context) {
	daysStr := c.Query("days")
	days, e := strconv.Atoi(daysStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 2000, e.Error())
	}
	res, err := discussionService.GetTrend(c, days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 1000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetDiscussionNewFeed(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 2000, e.Error())
		return

	}
	res, err := discussionService.GetNewFeed(c, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 1000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetDiscussionActiveUser(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, 2000, e.Error())
		return
	}
	res, err := discussionService.GetActiveUser(c, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 1000, err.Error())
		return
	}
	common.Success(c, res)
}
