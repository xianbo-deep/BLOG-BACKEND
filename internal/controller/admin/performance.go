package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAverageDelay(c *gin.Context) {
	res, err := admin.GetAverageDelay()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}

func GetSlowPages(c *gin.Context) {
	ctx := c.Request.Context()
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	res, err := admin.GetSlowPages(ctx, limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)
}
