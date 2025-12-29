package admin

import (
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, 1000, "invalid request")
		return
	}

	token, err := admin.AdminLogin(req.Username, req.Password)

	if err != nil {
		common.Fail(c, http.StatusUnauthorized, 1001, err.Error())
		return
	}

	common.Success(c, gin.H{
		"token":   token,
		"message": "Login Success",
	})
}
