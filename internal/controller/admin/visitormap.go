package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetVisitorMap(c *gin.Context) {
	// 传递毫秒级时间戳
	startTimeStr := c.DefaultQuery("startTime", "0")
	endTimeStr := c.DefaultQuery("endTime", "0")

	startTime, _ := strconv.Atoi(startTimeStr)
	endTime, _ := strconv.Atoi(endTimeStr)

	res, err := admin.GetVisitorMap(startTime, endTime)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)

}
