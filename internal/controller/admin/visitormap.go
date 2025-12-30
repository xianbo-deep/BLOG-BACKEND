package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVisitorMap(c *gin.Context) {
	res, err := admin.GetVisitorMap()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, 2000, err.Error())
		return
	}
	common.Success(c, res)

}
