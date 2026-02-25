package core

import (
	"Blog-Backend/consts"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/oschwald/geoip2-golang"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	RDB   *redis.Client
	GeoDB *geoip2.Reader
	once  sync.Once
	IP2R4 *xdb.Searcher
	IP2R6 *xdb.Searcher
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

		if err := initIP2Region(); err != nil {
			initErr = fmt.Errorf("initialize IP2Region: %w", err)
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

	/* 初始化日志打印器 */
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		},
	)
	var err error

	/* 获取实例 */
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
		Logger:      newLogger,
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

/* 初始化ip2region */
func initIP2Region() error {
	v4 := os.Getenv("IP2REGION_V4_PATH")
	v6 := os.Getenv("IP2REGION_V6_PATH")

	if v4 == "" || v6 == "" {
		return fmt.Errorf("%w", consts.ErrIP2RegionDBNotFound)
	}
	IP2R4 = loadSearcher(v4)
	IP2R6 = loadSearcher(v6)
	return nil
}

/* 关闭ip2region */
func closeIP2RegionDB() {
	if IP2R4 != nil {
		IP2R4.Close()
	}
	if IP2R6 != nil {
		IP2R6.Close()
	}
}

// 全量加载到内存
func loadSearcher(path string) *xdb.Searcher {
	// 从文件解析出version
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open xdb failed: %v", err)
	}
	defer f.Close()

	// 校验xdb
	if err := xdb.Verify(f); err != nil {
		log.Fatalf("verify xdb failed: %v", err)
	}

	// 获取头部
	header, err := xdb.LoadHeader(f)
	if err != nil {
		log.Fatalf("load xdb header failed: %v", err)
	}

	// 获取版本
	version, err := xdb.VersionFromHeader(header)
	if err != nil {
		log.Fatalf("detect xdb version failed: %v", err)
	}

	// 全量加载到内存
	buff, err := xdb.LoadContentFromFile(path)
	if err != nil {
		log.Fatalf("load xdb failed: %v", err)
	}

	searcher, err := xdb.NewWithBuffer(version, buff)
	if err != nil {
		log.Fatalf("load xdb failed: %v", err)
	}
	return searcher
}
