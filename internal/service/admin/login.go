package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/utils"
	"errors"
)

type LoginService struct {
}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (s *LoginService) AdminLogin(Username string, Password string) (string, error) {
	// TODO 后续加入数据库管理 不需要在环境变量存储
	adminUser := consts.EnvAdminUser
	adminPwd := consts.EnvAdminPwd

	if Username != adminUser || Password != adminPwd {
		return "", errors.New("invalid username or password")
	}

	// 生成token
	token, err := utils.GenerateToken(Username)

	// 判断是否出错
	if err != nil {
		return "", err
	}

	return token, nil
}
