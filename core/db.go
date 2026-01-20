package core

import (
	"Blog-Backend/consts"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	RDB   *redis.Client
	GeoDB *geoip2.Reader
	Ctx   = context.Background()
	once  sync.Once
)

func Init() error {
	var initErr error
	once.Do(func() {
		if err := initPG(); err != nil {
			initErr = fmt.Errorf("initialize PostgreSQL: %w", err)
			return
		}

		if err := initRedis(); err != nil {
			initErr = fmt.Errorf("initialize Redis: %w", err)
			return
		}

		if err := initGeoDB(); err != nil {
			initErr = fmt.Errorf("initialize GeoIP: %w", err)
			return
		}
	})
	return initErr
}

/* 初始化PG */
func initPG() error {
	/* 获取字符串 */
	dsn := os.Getenv(consts.EnvPgURI)

	if dsn == "" {
		return consts.ErrPostgresNotConfigured
	}

	var err error

	/* 获取实例 */
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})

	/* 判断是否出错 */
	if err != nil {
		return fmt.Errorf("%w: %v", consts.ErrDBConnectionFailed, err)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		return fmt.Errorf("get database instance: %w", err)
	}
	/* 配置连接池 */
	sqlDB.SetMaxOpenConns(25)           // 最大连接数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour) // 每个连接1h换1次
	return nil
}

/* 初始化Redis */
func initRedis() error {
	addr := os.Getenv(consts.EnvRedisURL)

	if addr == "" {
		return consts.ErrRedisNotConfigured
	}

	// 解析
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return fmt.Errorf("parse Redis URL: %w", err)
	}

	RDB = redis.NewClient(opt)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("%w: %v", consts.ErrRedisConnectionFail, err)
	}
	return nil
}

/* 初始化GeoDB */
func initGeoDB() error {
	var err error
	GeoDB, err = geoip2.Open(os.Getenv(consts.EnvGeoDBPath))
	if err != nil {
		return fmt.Errorf("%w: %v", consts.ErrGeoDBNotFound, err)
	}
	return nil
}
