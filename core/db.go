package core

import (
	"Blog-Backend/consts"
	"context"
	"log"
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

func Init() {
	once.Do(func() {
		initPG()
		initRedis()
		initGeoDB()
	})
}

/* 初始化PG */
func initPG() {
	/* 获取字符串 */
	dsn := os.Getenv(consts.EnvPgURI)

	if dsn == "" {
		log.Fatal("Lack of URL of pgsql")
	}

	var err error

	/* 获取实例 */
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})

	/* 判断是否出错 */
	if err != nil {
		log.Fatal(err)
		return
	}

	sqlDB, _ := DB.DB()

	/* 配置连接池 */
	sqlDB.SetMaxOpenConns(25)           // 最大连接数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour) // 每个连接1h换1次

}

/* 初始化Redis */
func initRedis() {
	addr := os.Getenv(consts.EnvRedisURL)

	if addr == "" {
		log.Fatal("Lack of URL of redis")
	}

	// 解析
	opt, err := redis.ParseURL(addr)
	if err != nil {
		log.Fatal("Invalid redis URL")
	}

	RDB = redis.NewClient(opt)

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatal("Fail Connection to redis")
	}
}

/* 初始化GeoDB */
func initGeoDB() {
	var err error
	GeoDB, err = geoip2.Open(consts.EnvGeoDBPath)
	if err != nil {
		log.Fatal("Failed to open GeoLite2-City.mmdb", err)
	}
}
