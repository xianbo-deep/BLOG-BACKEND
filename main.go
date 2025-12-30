package main

import (
	"Blog-Backend/core"
	"Blog-Backend/router"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	rEngine *gin.Engine
	once    sync.Once
)

func main() {
	/* 初始化数据库 */
	core.Init()

	/* 初始化路由 */
	once.Do(func() {
		rEngine = router.SetupRouter()
	})
	/* 获取端口 */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	/* 监听端口，启动后端 */
	if err := rEngine.Run(":" + port); err != nil {
		panic(err)
	}

}
