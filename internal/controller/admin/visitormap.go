package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var visitorMapService *admin.VisitorMapSerive

func InitVisitorMapService() {
	visitorMapService = admin.NewVisitorMapSerive()
}

func GetVisitorMap(c *gin.Context) {
	// 传递毫秒级时间戳
	startTimeStr := c.DefaultQuery("startTime", "0")
	endTimeStr := c.DefaultQuery("endTime", "0")

	startTime, _ := strconv.Atoi(startTimeStr)
	endTime, _ := strconv.Atoi(endTimeStr)

	res, err := visitorMapService.GetVisitorMap(startTime, endTime)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)

}

func GetChineseVisitorMap(c *gin.Context) {
	// 传递毫秒级时间戳
	startTimeStr := c.DefaultQuery("startTime", "0")
	endTimeStr := c.DefaultQuery("endTime", "0")

	startTime, _ := strconv.Atoi(startTimeStr)
	endTime, _ := strconv.Atoi(endTimeStr)

	res, err := visitorMapService.GetVisitorMap(startTime, endTime)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, res)
}
