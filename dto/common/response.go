package common

import (
	"Blog-Backend/consts"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    consts.CodeSuccess,
		Message: consts.ErrorMessage(consts.CodeSuccess),
		Data:    data,
	})
}

func Fail(c *gin.Context, httpCode int, code int, errMsg string) {
	message := consts.ErrorMessage(code)
	if message == "" {
		message = errMsg
	}

	response := Response{
		Code:    code,
		Message: message,
	}
	if errMsg != "" {
		response.Error = errMsg
	}
	c.JSON(httpCode, Response{
		Code:    response.Code,
		Message: response.Message,
		Error:   response.Error,
	})
}
