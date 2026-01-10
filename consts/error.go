package consts

import "fmt"

const (
	CodeSuccess      = 0
	CodeBadRequest   = 1000
	CodeUnauthorized = 1001
	CodeInternal     = 1005

	CodeInvalidToken    = 2000
	CodeTokenExpired    = 2001
	CodeInvalidPassword = 2002
	CodeUserNotFound    = 2003
	CodeTokenRequired   = 2005
)

var errorMessages = map[int]string{
	CodeSuccess:         "success",
	CodeBadRequest:      "invalid request",
	CodeUnauthorized:    "unauthorized",
	CodeInternal:        "internal server error",
	CodeInvalidToken:    "invalid token",
	CodeTokenExpired:    "token expired",
	CodeInvalidPassword: "invalid password",
	CodeUserNotFound:    "user not found",
	CodeTokenRequired:   "token is required",
}

func ErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return ""
}

// 通用的初始化错误
var (
	ErrPostgresNotConfigured = fmt.Errorf("PG database URL not configured")
	ErrRedisNotConfigured    = fmt.Errorf("Redis URL not configured")
	ErrGeoDBNotFound         = fmt.Errorf("GeoIP database file not found")
)

// 数据库相关错误
var (
	ErrDBConnectionFailed  = fmt.Errorf("database connection failed")
	ErrDBQueryFailed       = fmt.Errorf("database query failed")
	ErrDBMigrateFailed     = fmt.Errorf("database migration failed")
	ErrRedisConnectionFail = fmt.Errorf("Redis connection failed")
)
