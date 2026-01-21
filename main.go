package main

import (
	"Blog-Backend/bootstrap"
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/internal/task"
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
	err := core.Init()

	if err != nil {
		panic(err)
	}

	/* 初始化组件 它们依赖于数据库 */
	bootstrap.InitComponet()

	/* 启动定时任务*/
	task.InitCron()

	/* 初始化路由 */
	rEngine = router.SetupRouter()

	/* 获取端口 */
	port := os.Getenv(consts.EnvPort)
	if port == "" {
		port = "8080"
	}

	/* 监听端口，启动后端 */
	if err := rEngine.Run(":" + port); err != nil {
		panic(err)
	}

}
