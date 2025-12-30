package admin

import (
	"Blog-Backend/utils"
	"errors"
	"os"
)

func AdminLogin(Username string, Password string) (string, error) {
	adminUser := os.Getenv("ADMIN_USER")
	adminPwd := os.Getenv("ADMIN_PASSWORD")

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
