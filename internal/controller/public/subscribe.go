package public

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/service/public"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscribeController struct {
	svc *public.SubscribeService
}

func NewSubscribeController(svc *public.SubscribeService) *SubscribeController {
	return &SubscribeController{svc: svc}
}

func (ctrl *SubscribeController) SubscribeBlog(c *gin.Context) {
	email := c.Query("email")
	subscribeStr := c.Query("subscribe")
	vc := c.Query("vc")
	ctx := c.Request.Context()
	if email == "" || subscribeStr == "" || vc == "" {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, errors.New("请求参数错误").Error())
		return
	}
	subscribe, e := strconv.Atoi(subscribeStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}

	err := ctrl.svc.SubscribeBlog(ctx, email, vc, subscribe)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, consts.CodeSuccess)
}

func (ctrl *SubscribeController) VerifyEmail(c *gin.Context) {
	email := c.Query("email")
	subscribeStr := c.Query("subscribe")
	if email == "" {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, errors.New("缺失邮箱").Error())
		return
	}
	subscribe, e := strconv.Atoi(subscribeStr)
	if e != nil {
		common.Fail(c, http.StatusBadRequest, consts.CodeBadRequest, e.Error())
		return
	}
	ctx := c.Request.Context()
	err := ctrl.svc.VerifyEmail(ctx, email, subscribe)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	common.Success(c, consts.CodeSuccess)
}
