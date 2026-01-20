package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/utils"
	"errors"
	"os"
)

type LoginService struct {
}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (s *LoginService) AdminLogin(Username string, Password string) (string, error) {
	// TODO 后续加入数据库管理 不需要在环境变量存储
	adminUser := os.Getenv(consts.EnvAdminUser)
	adminPwd := os.Getenv(consts.EnvAdminPwd)

	if Username != adminUser {
		return "", errors.New(consts.ErrorMessage(consts.CodeUserNotFound))
	}
	if Password != adminPwd {
		return "", errors.New(consts.ErrorMessage(consts.CodeInvalidPassword))
	}
	// 生成token
	token, err := utils.GenerateToken(Username)

	// 判断是否出错
	if err != nil {
		return "", err
	}

	return token, nil
}
