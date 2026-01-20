package utils

import (
	"Blog-Backend/consts"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var Secret = []byte(os.Getenv(consts.EnvJWTSecret))

// 生成Token
func GenerateToken(username string) (string, error) {
	claims := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 24小时过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.JwtTokenExpireDuration)),
			Issuer:    consts.JwtIssuer,
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
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New(consts.ErrorMessage(consts.CodeInvalidToken))
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New(consts.ErrorMessage(consts.CodeTokenExpired))
		}
	}

	// 类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(consts.ErrorMessage(consts.CodeInvalidToken))
}
