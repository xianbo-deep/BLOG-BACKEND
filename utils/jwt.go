package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var Secret = []byte(os.Getenv("JWT_SECRET"))

// 生成Token
func GenerateToken(username string) (string, error) {
	claims := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "xbZhong",
		},
	}
	// HS256签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 密钥加密
	return token.SignedString(Secret)
}

// 解密token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})

	if err != nil {
		return nil, err
	}

	// 类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
