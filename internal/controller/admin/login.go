package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/dto/request"
	"Blog-Backend/internal/service/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	svc *admin.LoginService
}

func NewLoginController(svc *admin.LoginService) *LoginController {
	return &LoginController{svc: svc}
}

func (ctrl *LoginController) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, err.Error())
		return
	}

	token, err := ctrl.svc.AdminLogin(req.Username, req.Password)

	if err != nil {
		common.Fail(c, http.StatusUnauthorized, consts.CodeUnauthorized, err.Error())
		return
	}

	common.Success(c, gin.H{
		"token":   token,
		"message": "Login Success",
	})
}
